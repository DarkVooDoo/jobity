package handler

import (
	"fmt"
	"job/store"
	"log"
	"net/http"
)

var ProSignHandler = func(res http.ResponseWriter, req *http.Request){
    route, err := NewRoute(res, req)

    route.Get(func() {
        if err == nil{
            route.Response.Header().Add("Location", "/")
            route.Response.WriteHeader(http.StatusTemporaryRedirect)
            return
        }
        route.Render(nil, "route/protemplate.html", "route/prosign.html")
    })

    route.Post(nil, func() {
        password := route.UrlEncoded["password"]
        email := route.UrlEncoded["email"]

        user , err := store.SigninProUser(email, password)
        if err != nil{
            log.Println(err)
            route.Notification("error", "Mauvais email ou mot de passe")
            return
        }
        token, err :=  store.CreateTokenPro(user)
        if err != nil{
            log.Println(err)
            route.Notification("error", "Token error")
        }
        route.Response.Header().Add("Set-Cookie", fmt.Sprintf("x-auth=%v;path=/;SameSite=Strict;HttpOnly;Max-Age=%v",token, 60*60*6))
        route.Response.Header().Add("HX-Redirect", "/")
        route.Response.WriteHeader(http.StatusTemporaryRedirect)

    })
    
    route.Put(nil, func() {
        siret := route.UrlEncoded["siren"]
        email := route.UrlEncoded["email"]
        password := route.UrlEncoded["password"]
        confirmation := route.UrlEncoded["confirmation"]
        if password != confirmation || !store.IsValidPassword(password){
            log.Println("error confirmation and password doesnt match")
            route.Notification("error", "Mot de passe pas les memes")
            return
        }
        proUser := store.ProUser{Siren: siret, Email: email}
        if err := proUser.Create(password); err != nil{
            log.Println(err)
            route.Notification("error", "Impossible de creer cette compte")
            return
        }
        route.Response.Header().Add("HX-Redirect", "/")
        route.Response.WriteHeader(http.StatusTemporaryRedirect)
        log.Println(siret, email)
    })
    
    route.Delete(nil, func() {
        cookie := http.Cookie{
            Name: "x-auth",
            MaxAge: -1,
        }
        http.SetCookie(route.Response, &cookie)
        route.Response.Header().Add("HX-Redirect", "/connexion")
    })
}



