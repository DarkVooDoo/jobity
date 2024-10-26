package handler

import (
	"html/template"
	"job/store"
	"log"
	"net/http"
)

var ProParametre = func(res http.ResponseWriter, req *http.Request){
    route, err := NewRoute(res, req)
    if err != nil{
        log.Printf("Need to be connecred")
        return
    }

    pro := store.ProUser{Id: route.User.Id}
    route.Post(&pro, func() {
        conn := store.GetDBPoolConn()
        defer conn.Close()
        if err := pro.Modify(conn); err != nil{
            route.Notification("error", "error")
            return
        }
        temp, err := template.New("entreprise_name").Parse(`<h1 class="name" id="entreprise_name" hx-swap-oob="true">{{.}}</h1>`)
        if err != nil{
            log.Printf("error in the template: %v", err)
            return
        }
        if err := temp.Execute(route.Response, pro.Name); err != nil{
            log.Printf("error executing template: %v", err)
            return
        }
        route.Notification("success", "modifié")
    })
}
