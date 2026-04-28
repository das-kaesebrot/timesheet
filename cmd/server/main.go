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
	renderer, err := template.New("templates")
	if err != nil {
		log.Panicf("failed to load templates: %v", err)
	}

	h := handler.New(repo, renderer)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Root)

	mux.HandleFunc("GET /users", h.ListUsers)
	mux.HandleFunc("GET /users/new", h.NewUser)
	mux.HandleFunc("POST /users", h.CreateUser)
	mux.HandleFunc("GET /users/{id}", h.ShowUser)
	mux.HandleFunc("GET /users/{id}/edit", h.EditUser)
	mux.HandleFunc("POST /users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", h.DeleteUser)

	mux.HandleFunc("GET /users/{id}/entries", h.ListUserEntries)
	mux.HandleFunc("GET /users/{id}/entries/new", h.NewUserEntry)
	mux.HandleFunc("POST /users/{id}/entries", h.CreateUserEntry)

	mux.HandleFunc("GET /users/{id}/overview", h.OverviewUser)
	mux.HandleFunc("GET /users/{id}/export", h.ExportUser)

	mux.HandleFunc("GET /entries/{id}/edit", h.EditEntry)
	mux.HandleFunc("PATCH /entries/{id}", h.UpdateEntry)
	mux.HandleFunc("DELETE /entries/{id}", h.DeleteEntry)

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
