package handler

import (
	"job/store"
	"log"
	"net/http"
)

type ProfilePage struct{
    RequireData
    User store.User
    Curriculum store.Curriculum
}

var ProfileHandler = func(res http.ResponseWriter, req *http.Request){
    route, err := NewRoute(res, req)
    user := store.User{Id: route.User.Id}
    var curriculum store.Curriculum
    route.Get(func() {
        if err != nil{
            route.Response.Header().Add("Location", "/connexion")
            route.Response.WriteHeader(http.StatusTemporaryRedirect)
            return
        }
        user.GetProfile()
        curriculum.Get(store.EncryptCurriculumId(route.User.Id))
        data := ProfilePage{
            RequireData{Search: SearchQuery{Query: ""}},
            user,
            curriculum,
        }
        route.Render(data, "route/template.html", "route/profile.html")
    })

    route.Post(&user, func() {
        if err := user.Modify(); err != nil{
            log.Println(err)
            route.Notification("error", "error modification de utilisateur")
            return
        }
        route.Notification("success", "modification enregistrée")
    })

    route.Put(nil, func() {
        if route.Multipart.File.Size > 50000{
            return
        }
        
        log.Println(route.Multipart.File.Size)
        user.UploadPhoto(route.Multipart.File.Buffer)
    })

    route.Patch(nil, func() {
        log.Println(route.Multipart.File.Size)
    })
}
