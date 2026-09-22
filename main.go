package main

import (
	"embed"
	"log"

	"github.com/das-kaesebrot/timesheet/cmd/server"
)

var (
	Version = "v0.0.1-dev"
	GitHash = "0000000000000000000000000000000000000000"
)

//go:embed web
var webFS embed.FS

func main() {
	if err := server.Run(webFS, Version, GitHash); err != nil {
		log.Fatal(err)
	}
}
