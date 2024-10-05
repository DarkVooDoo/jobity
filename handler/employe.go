package handler

import (
	"html/template"
	"job/store"
	"log"
	"net/http"
)

var EmployeHandler = func(res http.ResponseWriter, req *http.Request){
    route, _ := NewRoute(res, req)
    
    var employe store.Employe = store.Employe{EntrepriseId: route.User.Id}
    route.Get(func() {
        employe.Email = route.Request.URL.Query().Get("email")
        if employe.Email == ""{
            return    
        }
        employeList := employe.Search()
        temp, err := template.New("Suggest Employe").Parse(`
            <div class="employe-search-suggest {{if eq (len .) 0}}hidden{{end}}">
                {{range .}}
                <form class="employe-suggest-preview" hx-post="/employe" hx-ext="json-enc" hx-vals='{"userId": "{{.UserId}}"}' 
                    hx-trigger="click" hx-swap="outerHTML" hx-target="closest .employe-search-suggest">
                        <div class="preview-photo">{{.ShortName}}</div>
                        <p>{{.Email}}</p>
                    </form>
                {{end}}
            </div>
        `)
        if err != nil{
            log.Printf("error creating template: %v", err)
        }
        if err := temp.Execute(route.Response, employeList); err != nil{
            log.Printf("error executing template: %v", err)
        }
    })

    route.Put(&employe, func() {
        if employe.ShopId != ""{
            if err := employe.ShopAssign(); err != nil{
                route.Notification("error", "requete error")
                route.Response.WriteHeader(http.StatusForbidden)
                return
            }
            route.Notification("success", "employé ajouté")
        }
    })

    route.Post(&employe, func() {
        if err := employe.New(); err != nil{
            log.Printf("error creating new employe")
            route.Notification("warning", "error")
            return
        }
        temp, _ := template.New("Reset Suggest Employe").Parse(`<div class="employe-search-suggest hidden"></div>`)
        route.Notification("success", "nouveau employé ajouté")
        temp.Execute(route.Response, nil)       
    })

    route.Delete(nil, func() {
        employe.Id = route.UrlEncoded["id"]
        if err := employe.Delete(); err != nil{
            route.Notification("error", "error")
            return
        }
        route.Notification("success", "employé supprimer")
    })
}
