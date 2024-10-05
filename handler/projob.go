package handler

import (
	"job/store"
	"log"
	"net/http"
)

type ProjobPage struct{
    RequireData
    Job store.Job
    Candidate []store.Curriculum
    Interview []store.Curriculum
    InterviewType []string
}

var ProJobHandler = func(res http.ResponseWriter, req *http.Request){
    route, err := NewRoute(res, req)

    var job store.Job = store.Job{Id: route.Request.PathValue("id"), EntrepriseId: route.User.Id}

    route.Get(func() {
        if err != nil{
            route.Response.Header().Add("Location", "/connexion")
            route.Response.WriteHeader(http.StatusTemporaryRedirect)
            return 
        }
        if err = job.GetJobById(); err != nil{
            log.Println(err)
            return
        }
        candidates, interview := store.GetJobCurriculum(job.Id)
        page := ProjobPage{
            RequireData{Search: SearchQuery{Query: ""}},
            job,
            candidates,
            interview,
            store.GetInterviewType(),
        }
        route.Render(page, "route/protemplate.html", "route/projob.html")
    })

    route.Post(nil, func() {
        if err := job.SaveAsTemplate(route.UrlEncoded["tname"]); err != nil{
            route.Notification("error", "error dans la requete")
            return
        }
        route.Notification("success", "modéle enregistrer")
    })

    route.Put(nil, func() {
        //var application = store.JobApplication{}       
        //application.UpdateStatus()
    })

    route.Patch(&job, func() {
        //Edit job offert
        if err := job.ModifyJob();err != nil{
            route.Notification("error", "error dans la modification")
            return
        }
        route.Notification("success", "offre modifié")
    })

    route.Delete(nil, func() {
        //Delete job application
        if err := job.DeleteJob(); err != nil{
            route.Notification("error", "Impossible de supprimer")
            return
        }
        route.Response.Header().Add("HX-Redirect", "/")
    })
}
