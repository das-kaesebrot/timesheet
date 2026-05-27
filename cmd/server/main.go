package main

import (
	"log"
	"net/http"
	"os"

	"github.com/das-kaesebrot/timesheet/internal/handler"
	"github.com/das-kaesebrot/timesheet/internal/middleware"
	"github.com/das-kaesebrot/timesheet/internal/model"
	"github.com/das-kaesebrot/timesheet/internal/repository"
	"github.com/das-kaesebrot/timesheet/internal/template"
	"github.com/glebarez/sqlite"
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

	var handlerMiddleware = []middleware.Middleware{
		middleware.LoggerMiddleware,
	}

	eh := middleware.ErrorHandler(renderer)
	with := func(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
		return middleware.Chain(eh(fn), handlerMiddleware...)
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	mux.HandleFunc("/", with(h.Root))

	mux.HandleFunc("GET /favicon.ico", eh(h.GetFavicon))

	mux.HandleFunc("GET /users", with(h.GetUsersList))
	mux.HandleFunc("GET /users/new", with(h.GetUserNew))
	mux.HandleFunc("POST /users", with(h.PostUserNew))
	mux.HandleFunc("GET /users/{id}", with(h.GetUserOverview))
	mux.HandleFunc("GET /users/{id}/edit", with(h.GetUserEdit))
	mux.HandleFunc("POST /users/{id}", with(h.PostUserUpdate))
	mux.HandleFunc("POST /users/{id}/delete", with(h.PostUserDelete))

	mux.HandleFunc("GET /users/{id}/entries/new", with(h.GetEntryNew))
	mux.HandleFunc("GET /users/{id}/entries/new/quick", with(h.GetEntryNewQuick))
	mux.HandleFunc("POST /users/{id}/entries", with(h.PostEntryNew))

	mux.HandleFunc("GET /users/{id}/export", with(h.ExportUser))

	mux.HandleFunc("GET /entries/{id}/edit", with(h.GetEntryEdit))
	mux.HandleFunc("POST /entries/{id}", with(h.PostEntryUpdate))
	mux.HandleFunc("POST /users/{id}/entries/delete", with(h.PostEntryDelete))

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
