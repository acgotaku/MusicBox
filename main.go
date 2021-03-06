package main

import (
	"MusicBox/service"
	// "MusicBox/service/qq"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Port string `json:"port"`
}

func main() {
	var path string
	flag.StringVar(&path, "c", "config/config.json", "Please use '-c' input config File")
	flag.Parse()
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var config Config
	json.Unmarshal(raw, &config)
	http.HandleFunc("/api/search", service.SearchHandler)
	http.HandleFunc("/api/track", service.TrackHandler)
	http.HandleFunc("/api/playlist", service.PlayListHandler)
	log.Printf("Start listen serve %s", config.Port)
	http.ListenAndServe(config.Port, nil)

}
