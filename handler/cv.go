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

    var cv store.Curriculum = store.Curriculum{UserId: route.Request.PathValue("userId"), JobId: route.Request.PathValue("id")}
    route.Get(func() {
        //Handle if the user isnt connected
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
        if err := cv.Get(cv.UserId); err != nil{
            log.Println(err)
            return
        }
        application := store.JobApplication{UserId: store.DecryptCurriculumId(route.Request.PathValue("userId")), JobId: route.Request.PathValue("id")}
        application.UpdateStatus("Vue")
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

    route.Delete(nil, func() {
        application := store.JobApplication{UserId: store.DecryptCurriculumId(route.Request.PathValue("userId")), JobId: route.Request.PathValue("id")}
        if err := application.UpdateStatus("Reject"); err != nil{
            route.Notification("error", "error requete impossible")
            return
        }
        route.Response.WriteHeader(http.StatusOK)
    })
}
