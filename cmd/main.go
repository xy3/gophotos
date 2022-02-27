package main

import (
	"encoding/json"
	"github.com/xy3/photos"
	"github.com/xy3/photos/handlers"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func loadConfig() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	configFile := path.Join(cwd, "config.json")
	log.Printf("Loading config from: %s\n", configFile)
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &photos.Config)
	if err != nil {
		return err
	}
	return nil
}

func writeConfig() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	configFile := path.Join(cwd, "config.json")
	log.Printf("Writing config to: %s\n", configFile)
	configJson, err := json.Marshal(photos.Config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, configJson, 0644)
}

func main() {
	err := loadConfig()
	if err != nil {
		log.Println("Failed to find or read a config.json file, using default config values instead")
		_ = writeConfig()
	}
	if err = photos.DbConnect(); err != nil {
		log.Fatal(err)
	}
	defer photos.DB.Close()

	if err = photos.DbSetup(photos.DB); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	handlers.SetupRoutes(mux)
	log.Printf("Photos Server Listening on %s%s\n", photos.Config.Host, photos.Config.BasePath)
	http.ListenAndServe(photos.Config.Host, mux)
}
