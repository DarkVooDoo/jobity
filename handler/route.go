package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"job/lib"
	"job/store"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type RouteInterface interface {
	Get()
	Post()
	Delete()
	Patch()
}

type SearchQuery struct{
    Query string
    City string
    Postal string
    Lat string
    Long string
    Lastposition int
    FrancetravailPosition int
}

type RequireData struct{
    User store.ConnectedUser
    Search SearchQuery
}

type FilePayload struct{
    Buffer io.Reader
    Size int
}

type Multipart struct{
    Body map[string]string
    File FilePayload
}

type ErrToken error

type Route struct {
	RouteInterface
    Multipart Multipart
	Params     map[string]string
    User store.ConnectedUser
	Request    *http.Request
	Response   http.ResponseWriter
	UrlEncoded map[string]string
    ContentType string
}

func NewRoute(response http.ResponseWriter, request *http.Request) (*Route, ErrToken) {

    var user store.ConnectedUser
    var multip Multipart = Multipart{Body: map[string]string{}}
	contentType := request.Header.Get("Content-Type")
    token, err := request.Cookie("x-auth")
    if err == nil{
        userData, accessToken, errToken := store.VerifyToken(token.Value)
        if errToken != nil{
            err = errToken
        }else{
            user = userData
            authCookie := http.Cookie{
                Name: "x-auth",
                Value: accessToken,
                MaxAge: 60*60*6,
                Path: "/",
                HttpOnly: true,
                SameSite: http.SameSiteStrictMode,
            }
            http.SetCookie(response, &authCookie)
        }
    }
    ct, params, _ := mime.ParseMediaType(request.Header.Get("Content-Type"))
	var urlEncoded map[string]string = map[string]string{}
	if contentType == "application/x-www-form-urlencoded" {
		urlEncoded = lib.ReadBody(request.Body)
    }else if ct == "multipart/form-data"{
        reader := multipart.NewReader(request.Body, params["boundary"])
        body := make([]byte, 1024)
        var payload string
        Parts:
        for{
            part, err := reader.NextPart()
            if err == io.EOF{
                break Parts
            }
            defer part.Close()
            var fileBuffer bytes.Buffer
            Exit:
            for {
                bytesReaded, err := part.Read(body)
                fileBuffer.Write(body[:bytesReaded])
                payload = string(body[:bytesReaded])
                if err == io.EOF {break Exit}
            }
            
            if part.FileName() != ""{
                multip.File.Buffer = &fileBuffer
                multip.File.Size = fileBuffer.Len()
                continue
            }
            multip.Body[part.FormName()] = payload
            
        }
    }	
    return &Route{
        ContentType: contentType,
        Multipart: multip,
		Request:    request,
		Response:   response,
        User: user,
		Params:     map[string]string{},
		UrlEncoded: urlEncoded,
	}, err
}

func DecodeJson(jsonStruct any, body io.ReadCloser){
    decoder := json.NewDecoder(body)
    err := decoder.Decode(jsonStruct)
    if err != nil{
        log.Println(err)
        log.Println("error decoding")
    }
}

func (r *Route) Post(jsonStruct any, handleFunc func()) {
	if r.Request.Method == http.MethodPost {
        if r.ContentType == "application/json" && jsonStruct != nil{
            DecodeJson(jsonStruct, r.Request.Body)  
        }
		handleFunc()
	}
}

func (r *Route) Patch(jsonStruct any, handleFunc func()) {
	if r.Request.Method == http.MethodPatch {
        if r.ContentType == "application/json" && jsonStruct != nil{
            DecodeJson(jsonStruct, r.Request.Body)  
        }
		handleFunc()
	}
}

func (r *Route) Put(jsonStruct any, handleFunc func()) {
	if r.Request.Method == http.MethodPut {
        if r.ContentType == "application/json" && jsonStruct != nil{
            DecodeJson(jsonStruct, r.Request.Body)  
        }
		handleFunc()
	}
}

func (r *Route) Delete(jsonStruct any, handleFunc func()) {
    if r.Request.Method == http.MethodDelete {
        if r.ContentType == "application/json" && jsonStruct != nil{
            DecodeJson(jsonStruct, r.Request.Body)  
        }

        handleFunc()
	}
}

func (r *Route) Get(handlerFunc func()) {
	if r.Request.Method == http.MethodGet {
		handlerFunc()
    }
}

func (r *Route) Render(payload interface{}, files ...string){
    temp, err := template.ParseFiles(files...)
    if err != nil{
        log.Println(err)
        return 
    }
    style, _ := compressCSS(files...)
    temp.New("Style").Parse(style)
    //temp.New("Body").Parse(body)
    r.Response.Header().Add("Content-Encoding", "gzip")
    gzip := gzip.NewWriter(r.Response)
    defer gzip.Close()
    if err := temp.Execute(gzip, payload); err != nil{
        log.Println(err)
        return
    }
}

func (r *Route) Notification(t string, msg string){
    notif := fmt.Sprintf(`<notification-ele id="notif" hx-swap-oob="true" msg="%v" type="%v" on></notification-ele>`, msg, t)
    templ, _ := template.New("notification").Parse(notif)
    if err := templ.Execute(r.Response, nil); err != nil{
        log.Println(err)
    }
}

func (r *Route) Redirect(path string){
    r.Response.Header().Add("Location", path)
    r.Response.WriteHeader(http.StatusTemporaryRedirect)
}

func compressCSS(files ...string)(styles string, body string){

    stylesStartCompress := `{{define "Style"}}`
    bodyStartCompress := `{{define "Body"}}`
    var filePath string
    for i := 0; i < len(files); i++{
        if strings.HasSuffix(files[i], "html"){
            filePath = files[i]
        }
    }
    file, err := os.Open(filePath)
    if err != nil{
        log.Println(err)
        return "", ""
    }
    defer file.Close()
    buf := make([]byte, 1024)
    result := ""
    var startPos, endPos, startBodyPos, endBodyPos int
    Read:
    for{
        readed, err := file.Read(buf)
        if err == io.EOF{
            break Read
        } 
        result += string(buf[:readed])
        startPos = strings.Index(result, stylesStartCompress)
        endPos = strings.Index(result[startPos:], "%end")
        if startPos > 0 && endPos > startPos{
            styles = strings.ReplaceAll(strings.ReplaceAll(result[startPos+len(stylesStartCompress):endPos+startPos-8], "\n", ""), "  ", "")
        }
        startBodyPos = strings.Index(result, bodyStartCompress) 
        if startBodyPos >= 0{
            endBodyPos = strings.Index(result[startBodyPos:], "%end")
            if endBodyPos > startBodyPos{
                body = strings.ReplaceAll(strings.ReplaceAll(result[startBodyPos+len(bodyStartCompress):endBodyPos+startBodyPos-8], "\n", ""), "  ", "")
            }
        }
    }
    return 
}
