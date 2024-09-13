package handler

import (
	"job/store"
	"log"
	"net/http"
)

type CurriculumPage struct{
    Curriculum store.Curriculum
}


var CurriculumHandler = func(res http.ResponseWriter, req *http.Request){
    route, _ := NewRoute(res, req)

    var cv store.Curriculum
    route.Get(func() {
        //Download Curriculum format pdf
        format := route.Request.URL.Query().Get("type")
        if format == "pdf"{
            cv.Get(store.EncryptCurriculumId(route.User.Id))
            pdf := store.CreateCurriculumPDF(cv)
            route.Response.Header().Add("Content-Disposition", `attachment; filename="filename.pdf"`)
            if err := pdf.Output(res); err != nil{
                log.Println(err)
            }
            pdf.Close()
            return
        }
        userId := route.Request.PathValue("userId")
        if err := cv.Get(userId); err != nil{
            log.Println(err)
            return
        }
        pageData := CurriculumPage{Curriculum: cv}
        route.Render(pageData, "route/protemplate.html", "route/curriculum.html")
    })

    route.Post(&cv, func() {
        //handle if the user isnt connected
        if err := cv.Save(route.User.Id); err != nil{
            log.Println(err)
            route.Notification("error", "error dans la requete")
            return
        }
        route.Notification("success", "Curriculum sauvegarder")

    })
}
