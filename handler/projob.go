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
    Possible []store.Curriculum
    InterviewType []string
    Category []store.Category
    Subcategory []store.Category
    Contract []store.Contract
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
        conn := store.GetDBPoolConn()
        defer conn.Close()
        if err := job.GetJobById(conn); err != nil{
            log.Printf("error getting the job")
            return
        }
        candidates, possible := store.GetJobCurriculum(conn, job.Id)
        page := ProjobPage{
            RequireData{Search: SearchQuery{Query: ""}},
            job,
            candidates,
            possible,
            store.GetInterviewType(conn),
            store.GetCategorys(conn),
            store.GetSubcategory(conn, job.CategoryId),
            store.GetContracts(conn),
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
        conn := store.GetDBPoolConn()
        defer conn.Close()
        if err := job.DeleteJob(conn); err != nil{
            route.Notification("error", "Impossible de supprimer")
            return
        }
        route.Response.Header().Add("HX-Redirect", "/")
    })
}
