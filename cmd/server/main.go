package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/das-kaesebrot/timesheet/internal/handler"
	"github.com/das-kaesebrot/timesheet/internal/middleware"
	"github.com/das-kaesebrot/timesheet/internal/model"
	"github.com/das-kaesebrot/timesheet/internal/repository"
	"github.com/das-kaesebrot/timesheet/internal/template"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var Version = "dev"

// returns nil if file exists and is read/writable, otherwise returns the underlying err
func checkFileAccess(path string) error {
	// doesnt't exist -> not read/writeable
	if _, err := os.Stat(path); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func main() {
	versionFlag := false
	flag.BoolVar(&versionFlag, "v", false, "print version information")
	flag.BoolVar(&versionFlag, "version", false, "print version information")
	flag.Parse()

	if versionFlag {
		fmt.Printf("%v\n", Version)
		return
	}

	log.Printf("Version: %v", Version)

	dbFile := path.Clean(os.Getenv("TIMESHEET_DB_FILE"))
	if dbFile == "." {
		dbFile = "timesheet.db"
	}

	log.Printf("Reading SQLite database from '%s'", dbFile)
	err := checkFileAccess(dbFile)

	if errors.Is(err, os.ErrNotExist) {
		log.Printf("database file doesn't exist yet, creating it")
	} else if err != nil {
		log.Panicf("failed reading database file: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.TimesheetEntry{})

	webDir := path.Clean(os.Getenv("TIMESHEET_WEB_DIR"))

	if webDir == "." {
		webDir = "web"
	}

	repo := repository.New(db)
	renderer, err := template.New(path.Join(webDir, "template"), Version)
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

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path.Join(webDir, "static")))))
	mux.HandleFunc("/", with(h.Root))

	mux.HandleFunc("GET /favicon.ico", eh(h.GetFavicon))

	mux.HandleFunc("GET /users", with(h.GetUsersList))
	mux.HandleFunc("GET /users/new", with(h.GetUserNew))
	mux.HandleFunc("POST /users", with(h.PostUserNew))
	mux.HandleFunc("GET /users/{id}", with(h.GetUserOverview))
	mux.HandleFunc("GET /users/{id}/edit", with(h.GetUserEdit))
	mux.HandleFunc("POST /users/{id}", with(h.PostUserUpdate))
	mux.HandleFunc("POST /users/{id}/delete", with(h.PostUserDelete))

	mux.HandleFunc("GET /users/{id}/entries", with(h.GetEntryNew))
	mux.HandleFunc("GET /users/{id}/entries/quick", with(h.GetEntryNewQuick))
	mux.HandleFunc("POST /users/{id}/entries", with(h.PostEntryNew))

	mux.HandleFunc("GET /users/{id}/entries/export", with(h.ExportUser))
	mux.HandleFunc("GET /users/{id}/entries/import", with(h.GetImportEntries))
	mux.HandleFunc("POST /users/{id}/entries/import", with(h.ImportEntriesToUser))

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
	log.Printf("Using '%s' as web dir", webDir)
	log.Fatal(http.ListenAndServe(host+":"+port, mux))
}
