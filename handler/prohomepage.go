package handler

import (
	"encoding/json"
	"fmt"
	"job/store"
	"log"
	"net/http"
)

type HomepageData struct{
    User store.ConnectedUser
    Profile store.ProUser
    Job []store.Job
}

type Adresse struct{
    Data []struct{
        Adresse struct{
            Name string `json:"name"`
            City string `json:"city"`
            Postal string `json:"postcode"`
        } `json:"properties"`
    } `json:"features"`
}

type Payload struct{
    Query string `jsonn:"city"`
}

var ProHomepageHandler = func(res http.ResponseWriter, req *http.Request){
    route, err := NewRoute(res, req)
    if req.URL.Path != "/"{
        route.Render("route/protemplate.htmk", "route/notfound.html")
        return
    }

    route.Get(func() {
        if err != nil{
            route.Response.Header().Add("Location", "/connexion")
            route.Response.WriteHeader(http.StatusTemporaryRedirect)
            return
        }
        q := route.Request.URL.Query().Get("city")
        if q != ""{
            //Search adresse from BAN api
            var addr Adresse
            req, _ := http.NewRequest("GET", fmt.Sprintf("https://api-adresse.data.gouv.fr/search/?q=%v", q), nil)
            res, err := http.DefaultClient.Do(req)
            if err != nil{
                log.Println(err)
            }
            dec := json.NewDecoder(res.Body)
            if err := dec.Decode(&addr); err != nil{
                log.Println(err)
            }
            log.Println(addr)
            return
        }
        var pageData HomepageData
        entreprise := store.ProUser{Id: route.User.Id}
        entreprise.GetProfile()
        jobs := store.GetEntrepriseJobs(route.User.Id)
        pageData.Profile = entreprise
        pageData.Job = jobs
        pageData.User = route.User
        route.Render(pageData, "route/protemplate.html", "route/prohomepage.html")
    })
}


