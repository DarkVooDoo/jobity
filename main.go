package main

import (
	"job/handler"
	"job/store"
	"log"
	"net/http"
)

const (
	Addr = ":8000"
)

func main() {
    if err := store.InitDB(); err != nil{
        log.Fatalf("error db init\n %v", err)
    }
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    Addr,
		Handler: mux,
	}
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("localhost/static/", http.StripPrefix("/static/", fs))
	mux.Handle("pro.localhost/static/", http.StripPrefix("/static/", fs))
    mux.HandleFunc("pro.localhost/", handler.ProHomepageHandler)
    mux.HandleFunc("pro.localhost/job/creer", handler.CreateJobHandler)
    mux.HandleFunc("pro.localhost/job/{id}", handler.ProJobHandler)
    mux.HandleFunc("pro.localhost/job/{id}/cv/{userId}", handler.CurriculumHandler)
    mux.HandleFunc("pro.localhost/template/{id}", handler.TemplateHandler)
    mux.HandleFunc("pro.localhost/connexion", handler.ProSignHandler)
    mux.HandleFunc("localhost/connexion", handler.SignHandler)
    mux.HandleFunc("localhost/job/{id}", handler.JobHandler)
    mux.HandleFunc("localhost/mes-postulations", handler.ApplicationHandler)
    mux.HandleFunc("localhost/profile", handler.ProfileHandler)
    mux.HandleFunc("localhost/profile/cv", handler.CurriculumHandler)
    mux.HandleFunc("localhost/entreprise/{id}", handler.CompanyHandler)
    mux.HandleFunc("localhost/search", handler.SearcHandler)
	mux.HandleFunc("localhost/", handler.HomepageRoute)
	if err := server.ListenAndServe(); err != nil {
        log.Println(err)
		log.Fatal("Keke crash")
	}
}
