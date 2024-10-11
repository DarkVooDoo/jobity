package handler

import (
	"job/store"
	"log"
	"net/http"
)

type CompanyPage struct{
    RequireData
    Job []store.Job
}

var CompanyHandler = func(res http.ResponseWriter, req *http.Request){
    route, _ := NewRoute(res, req)

    id := route.Request.PathValue("id")
    route.Get(func() {
        jobs := store.GetEntrepriseJobCards(id)
        log.Println(jobs)
        page := CompanyPage{
            RequireData{User: route.User},
            jobs,
        }
        route.Render(page, "route/template.html", "route/company.html")
    })
}
