package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xy3/photos"
	"github.com/xy3/photos/handlers"
	"log"
	"net/http"
	"strings"
)

var (
	host     = "localhost:8090"
	basePath = "/api/v1"
)

func main() {
	if err := photos.DbConnect(); err != nil {
		log.Fatal(err)
	}
	defer photos.DB.Close()

	fmt.Print("Do you want to run the database setup script? Y/n: ")
	var answer string
	_, _ = fmt.Scanln(&answer)

	if strings.ToLower(strings.TrimSpace(answer)) == "y" {
		if err := photos.DbSetup(photos.DB); err != nil {
			log.Fatal(err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc(basePath+"/users/signup", handlers.Signup)
	mux.HandleFunc(basePath+"/photos/download", photos.BasicAuthMiddleware(handlers.DownloadPhoto))
	mux.HandleFunc(basePath+"/photos", photos.BasicAuthMiddleware(handlers.PhotosHandler))
	mux.HandleFunc(basePath+"/users", photos.BasicAuthMiddleware(handlers.UsersHandler))

	log.Printf("Photos Server Listening on %s%s\n", host, basePath)
	http.ListenAndServe(host, mux)
}
