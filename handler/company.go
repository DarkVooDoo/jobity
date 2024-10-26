package handler

import (
	"job/store"
	"net/http"
)

type CompanyPage struct{
    RequireData
    Name string
    Adresse string
    Job []store.Job
}

var CompanyHandler = func(res http.ResponseWriter, req *http.Request){
    route, _ := NewRoute(res, req)

    id := route.Request.PathValue("id")
    route.Get(func() {
        conn := store.GetDBPoolConn()
        defer conn.Close()
        name, addr := store.GetEntrepriseInfo(conn, id)
        jobs := store.GetEntrepriseJobCards(conn, id)
        page := CompanyPage{
            RequireData{User: route.User},
            name,
            addr,
            jobs,
        }
        route.Render(page, "route/template.html", "route/company.html")
    })
}
