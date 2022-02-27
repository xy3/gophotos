package main

import (
	"github.com/xy3/photos"
	"github.com/xy3/photos/handlers"
	"log"
	"net/http"
)

func main() {
	if err := photos.DbConnect(); err != nil {
		log.Fatal(err)
	}
	defer photos.DB.Close()

	if err := photos.DbSetup(photos.DB); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	handlers.SetupRoutes(mux)
	log.Printf("Photos Server Listening on %s%s\n", photos.Config.Host, photos.Config.BasePath)
	http.ListenAndServe(photos.Config.Host, mux)
}
