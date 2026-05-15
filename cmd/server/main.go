package main

import (
	"log"
	"net/http"
	"os"

	"github.com/das-kaesebrot/timesheet/internal/handler"
	"github.com/das-kaesebrot/timesheet/internal/model"
	"github.com/das-kaesebrot/timesheet/internal/repository"
	"github.com/das-kaesebrot/timesheet/internal/template"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("timesheet.db"), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.TimesheetEntry{})

	repo := repository.New(db)
	renderer, err := template.New("web/template")
	if err != nil {
		log.Panicf("failed to load templates: %v", err)
	}

	h := handler.New(repo, renderer)

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	mux.HandleFunc("/", h.Root)

	mux.HandleFunc("GET /users", h.GetUsersList)
	mux.HandleFunc("GET /users/new", h.GetUserNew)
	mux.HandleFunc("POST /users", h.PostUserNew)
	mux.HandleFunc("GET /users/{id}", h.GetUserOverview)
	mux.HandleFunc("GET /users/{id}/edit", h.GetUserEdit)
	mux.HandleFunc("POST /users/{id}", h.PostUserUpdate)
	mux.HandleFunc("DELETE /users/{id}", h.PostUserDelete)

	mux.HandleFunc("GET /users/{id}/entries/new", h.GetEntryNew)
	mux.HandleFunc("GET /users/{id}/entries/new/quick", h.GetEntryNewQuick)
	mux.HandleFunc("POST /users/{id}/entries", h.PostEntryNew)
	mux.HandleFunc("POST /users/{id}/entries/quick", h.PostEntryNewQuick)

	mux.HandleFunc("GET /users/{id}/export", h.ExportUser)

	mux.HandleFunc("GET /entries/{id}/edit", h.GetEntryEdit)
	mux.HandleFunc("POST /entries/{id}", h.PostEntryUpdate)
	mux.HandleFunc("GET /entries/{id}/delete", h.PostEntryDelete)
	mux.HandleFunc("DELETE /entries/{id}", h.PostEntryDelete)

	host := os.Getenv("TIMESHEET_HOST")
	if host == "" {
		host = "[::]"
	}

	port := os.Getenv("TIMESHEET_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on host %s:%s", host, port)
	log.Fatal(http.ListenAndServe(host+":"+port, mux))
}
