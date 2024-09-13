package handler

import (
	"log"
	"net/http"
)


var TemplateHandler = func(res http.ResponseWriter, req *http.Request){
    route, _ := NewRoute(res, req)

    id := route.Request.PathValue("id")
    route.Post(nil, func() {
        name := route.UrlEncoded["name"]
        log.Println(id, name)
        route.Notification("success", "Modele enregistrer")
    })

}
