package handler

import "net/http"

var CompanyHandler = func(res http.ResponseWriter, req *http.Request){
    route, _ := NewRoute(res, req)

    route.Get(func() {
        route.Render(nil, "route/template.html", "route/company.html")
    })
}
