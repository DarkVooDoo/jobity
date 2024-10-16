package main

import (
	"job/handler"
	"job/store"
	"log"
	"net/http"
)

const (
	Addr = ":8000"
    Subdomaine = ":8001"
)

func Subdomain () {
	    mux := http.NewServeMux()
	    server := &http.Server{
	    	Addr:    Subdomaine,
	    	Handler: mux,
	    }
        fs := http.FileServer(http.Dir("./static"))
        mux.Handle("/static/", http.StripPrefix("/static/", fs))
        mux.HandleFunc("/", handler.ProHomepageHandler)
        mux.HandleFunc("/job/creer", handler.CreateJobHandler)
        mux.HandleFunc("/job/{id}", handler.ProJobHandler)
        mux.HandleFunc("/job/{id}/cv/{userId}", handler.CurriculumHandler)
        mux.HandleFunc("/boutiques", handler.ShopHandler)
        mux.HandleFunc("/employe", handler.EmployeHandler)
        mux.HandleFunc("/template/{id}", handler.TemplateHandler)
        mux.HandleFunc("/connexion", handler.ProSignHandler)
        if err := server.ListenAndServe(); err != nil {
            log.Println(err)
            log.Fatal("Keke crash")
        }
    }

func Maindomain(){
	    mux := http.NewServeMux()
	    server := &http.Server{
	    	Addr:    Addr,
	    	Handler: mux,
	    }
        fs := http.FileServer(http.Dir("./static"))
        mux.Handle("/static/", http.StripPrefix("/static/", fs))
        mux.HandleFunc("/connexion", handler.SignHandler)
        mux.HandleFunc("/job/{id}", handler.JobHandler)
        mux.HandleFunc("/mes-postulations", handler.ApplicationHandler)
        mux.HandleFunc("/bookmark", handler.BookmarkHandler)
        mux.HandleFunc("/profile", handler.ProfileHandler)
        mux.HandleFunc("/profile/cv", handler.CurriculumHandler)
        mux.HandleFunc("/entreprise/{id}", handler.CompanyHandler)
        mux.HandleFunc("/search", handler.SearcHandler)
        mux.HandleFunc("/", handler.HomepageRoute)
        if err := server.ListenAndServe(); err != nil {
            log.Println(err)
            log.Fatal("Keke crash")
        }
}

func main() {
    if err := store.InitDB(); err != nil{
        log.Fatalf("error db init\n %v", err)
    }
	finish := make(chan bool)
    go Maindomain()
    go Subdomain()
    <-finish
}
