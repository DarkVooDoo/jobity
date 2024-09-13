package main

import (
	"job/handler"
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
        mux.HandleFunc("/ile", handler.ProfileHandler)
        mux.HandleFunc("/ile/cv", handler.CurriculumHandler)
        mux.HandleFunc("/entreprise/{id}", handler.CompanyHandler)
        mux.HandleFunc("/search", handler.SearcHandler)
        mux.HandleFunc("/", handler.HomepageRoute)
        if err := server.ListenAndServe(); err != nil {
            log.Println(err)
            log.Fatal("Keke crash")
        }
}

func main() {
	finish := make(chan bool)
    go Maindomain()
    go Subdomain()
    <-finish
}
