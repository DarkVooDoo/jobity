package handler

import (
	"job/store"
	"log"
	"net/http"
)

type SignPage struct{
    RequireData
}


var SignHandler = func(res http.ResponseWriter, req *http.Request){
    var user store.User
    route, err := NewRoute(res, req)

    route.Get(func() {
        if err == nil{
            route.Response.Header().Add("Location", "/")
            route.Response.WriteHeader(http.StatusTemporaryRedirect )
            return
        }
        data := SignPage{RequireData{Search: SearchQuery{Query: ""}}}
        route.Render(data, "route/template.html", "route/sign.html")
    })

    route.Post(nil, func() {
        //Sign user
        var token string
        googleCredential := route.UrlEncoded["credential"]
        if googleCredential != ""{
            _, googleUser, err := store.VerifyGoogleToken(googleCredential)
            if err != nil{
                return
            }
            if err := store.SignGoogleUser(googleUser, googleCredential); err != nil{
                return
            }
            token = googleCredential
            authToken := http.Cookie{
                Name: "x-auth",
                Path: "/",
                SameSite: http.SameSiteStrictMode,
                HttpOnly: true,
                MaxAge: 60*60*6,
                Value: token,
            }
            http.SetCookie(route.Response, &authToken)
            route.Response.Header().Add("Location", "/")
            route.Response.WriteHeader(http.StatusMovedPermanently)
            return
            //Log in with google
        }else{
            email := route.UrlEncoded["email"]
            password := route.UrlEncoded["password"]
            user, err := store.SigninUser(email, password)
            if err != nil{
                log.Println(err)
                route.Notification("error", "Mauvais email ou mot de passe")
                return
            }
            token, _ = store.CreateToken(user)

            authToken := http.Cookie{
                Name: "x-auth",
                Path: "/",
                SameSite: http.SameSiteStrictMode,
                HttpOnly: true,
                MaxAge: 60*60*6,
                Value: token,
            }
            http.SetCookie(route.Response, &authToken)
            route.Response.Header().Add("HX-Redirect", "/")
            route.Response.WriteHeader(http.StatusTemporaryRedirect)
        }
    })
    
    route.Put(nil, func() {
        //Create User
        password := route.UrlEncoded["password"]
        confirmation := route.UrlEncoded["confirmation"]
        user.Email = route.UrlEncoded["email"]
        if confirmation != password && !store.IsValidPassword(password){
            //invalid password
            route.Notification("error", "Les mots de passe ne sont pas les mêmes")
            return
        }
        if err := user.Create(password); err != nil{
            log.Println(err)
            return
        }
    })

    route.Delete(nil, func() {
        deleteCookie := http.Cookie{
            Name: "x-auth",
            MaxAge: -1,
        }
        http.SetCookie(route.Response, &deleteCookie)
        route.Response.Header().Add("HX-Redirect", "/")
    })
}
