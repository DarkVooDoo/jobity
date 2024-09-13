package handler

import (
	"encoding/json"
	"html/template"
	"job/lib"
	"job/store"
	"log"
	"net/http"
	"strings"
)

type PagePayload struct{
    RequireData
    Job []store.Job
    FTJob []store.Job
    FTRecomendation []store.Job
}

var HomepageRoute = func(res http.ResponseWriter, req *http.Request) {
    data := PagePayload{RequireData: RequireData{Search: SearchQuery{Query: ""}}}
	if req.URL.Path != "/" {
		temp, _ := template.ParseFiles("route/template.html", "route/notfound.html")
        if err := temp.Execute(res, data); err != nil{
            log.Println(err)
            return
        }
        return
    }
    route, _ := NewRoute(res, req)
	route.Get(func() {
        var recomendation []lib.RecomendationToken
        var myRecomendation lib.RecomendationToken
        cookie, err := route.Request.Cookie("foryou")
        if err == nil{
            if err = json.Unmarshal([]byte(strings.ReplaceAll(cookie.Value, "%22", `"`)), &recomendation); err != nil{
                log.Println(err)
            }
            myRecomendation = lib.MostSearch(recomendation)
        }
        ftJobs, ftRecomendationJobs, jobs := store.GetEmplois(myRecomendation)
        data.User = route.User
        data.Job = jobs
        data.FTJob = ftJobs
        data.FTRecomendation = ftRecomendationJobs
        route.Render(data, "route/template.html", "route/page.html")
	})
}
