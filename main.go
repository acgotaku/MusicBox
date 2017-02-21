package main

import (
	"MusicBox/service"
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
	http.HandleFunc("/api/search", service.NeteaseSearchHandler)
	key := []byte("5d02ee4df79bb64b")
	plaintext := []byte("sigbgoZY+5dHJl5fY4Ri/HqkGvfWKHbAX+p1OK2wwYDqcy+1sfkBm0TliBddazu1")
	service.AesEncrypt(plaintext, key)
	log.Printf("Start listen serve %s", config.Port)
	http.ListenAndServe(config.Port, nil)

}
