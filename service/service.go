package service

import (
	"MusicBox/model"
	"MusicBox/service/netease"
	"MusicBox/service/qq"
	"MusicBox/service/xiami"
	"encoding/json"

	"net/http"
	"strconv"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	keyword := r.URL.Query().Get("keyword")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	source := r.URL.Query().Get("source")
	if limit == 0 {
		limit = 20
	}
	if page == 0 {
		page = 1
	}
	var musicDetail []model.MusicDetail
	switch source {
	case "qq":
		musicDetail = qq.SearchMusic(keyword, limit, page)
	case "xiami":
		musicDetail = xiami.SearchMusic(keyword, limit, page)
	default:
		musicDetail = netease.SearchMusic(keyword, limit, page)
	}
	music, _ := json.Marshal(musicDetail)
	w.Write(music)
}

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	track := r.URL.Query().Get("id")
	source := r.URL.Query().Get("source")

	var trackDetail model.TrackDetail
	switch source {
	case "qq":
		trackDetail = qq.GetTrack(track)
	case "xiami":
		trackDetail = xiami.GetTrack(track)
	default:
		trackDetail = netease.GetTrack(track)
	}
	trackJson, _ := json.Marshal(trackDetail)
	w.Write([]byte(trackJson))
}
