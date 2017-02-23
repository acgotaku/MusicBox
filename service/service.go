package service

import (
	"MusicBox/service/netease"
)

type MusicDetail struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	keyword := r.URL.Query().Get("keyword")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	source := r.URL.Query().Get("source")

	switch source {
	case "qq":
		qq.SearchMusic(keyword, limit, offset)
	case "xiami":
		xiami.SearchMusic(keyword, limit, offset)
	default:
		netease.SearchMusic(keyword, limit, offset)
	}
	w.Write(music)
}
