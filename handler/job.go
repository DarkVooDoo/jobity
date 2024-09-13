package handler

import (
	"encoding/json"
	"job/lib"
	"job/store"
	"log"
	"net/http"
	"strings"
)

type JobPage struct{
    RequireData
    store.Job 
}

var JobHandler = func(res http.ResponseWriter, req *http.Request){
    route, _ := NewRoute(res, req)

    id := route.Request.PathValue("id")
    job := store.Job{Id: id, EntrepriseId: route.User.Id}
    route.Get(func() {
        isThirdParty := route.Request.URL.Query().Get("tp")
        if isThirdParty == "true"{
            //Fetch france travail offert
            ftJob, err := store.GetFranceTravailJobById(id)
            if err != nil{
                log.Printf("error france travail: %V")
                route.Response.Header().Add("Location", "/")
                route.Response.WriteHeader(http.StatusTemporaryRedirect)
                return
            }
            job = ftJob
        }else{
            if err := job.GetJobById(); err != nil{
                log.Println(err)
                return
            }
        }
        page := JobPage{RequireData{Search: SearchQuery{Query: ""}, User: route.User}, job}
        var recomendation []lib.RecomendationToken
        cookie, err := req.Cookie("foryou")
        if err != nil{
            log.Printf("error cookie: %v", err)
        }else{
            foryou := strings.ReplaceAll(cookie.Value, "%22", `"`)
            if err = json.Unmarshal([]byte(foryou), &recomendation); err != nil{
                log.Println(err)
            }
        }

        //TODO: How to handle both our app and france travail api contract names
        //contract := lib.Contract[job.Contract]
        //log.Println(contract)
        label := strings.Fields(job.Title)[0]
        queue := lib.Queue[lib.RecomendationToken]{List: recomendation, Length: 5}
        if queue.IsFull(){
            queue.Dequeue()
            queue.Enqueue(lib.RecomendationToken{Postal: job.Postal[:2], Contract: job.Contract, Fulltime: job.Fulltime, Label: label})
        }else{
            queue.Enqueue(lib.RecomendationToken{Postal: job.Postal[:2], Contract: job.Contract, Fulltime: job.Fulltime, Label: label})
        }
        decRecomendation, _ := json.Marshal(queue.List)
        cookieStruct := http.Cookie{
            Name: "foryou",
            MaxAge: 60*60*24*5,
            Path: "/",
            SameSite: http.SameSiteStrictMode,
            HttpOnly: true,
            Value: strings.ReplaceAll(string(decRecomendation), `"`, "%22"),
        }
        http.SetCookie(route.Response, &cookieStruct)
        route.Render(page, "route/template.html", "route/job.html")

    })
    
}
