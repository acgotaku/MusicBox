package service

import (
	"MusicBox/model"
	"MusicBox/service/netease"
	"MusicBox/service/qq"
	"MusicBox/service/xiami"
	"encoding/json"
	"fmt"
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
	if keyword == "" {
		nullResponse := fmt.Sprintf(`{"code":404,"data":[],"total":0}`)
		w.Write([]byte(nullResponse))
		return
	}
	var musicSearch model.MusicSearch
	switch source {
	case "qq":
		musicSearch = qq.SearchMusic(keyword, limit, page)
	case "xiami":
		musicSearch = xiami.SearchMusic(keyword, limit, page)
	default:
		musicSearch = netease.SearchMusic(keyword, limit, page)
	}
	music, _ := json.Marshal(musicSearch)
	w.Write(music)
}

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	track := r.URL.Query().Get("id")
	source := r.URL.Query().Get("source")
	country := r.URL.Query().Get("country")
	if country == "" {
		country = "china"
	}
	if track == "" {
		nullResponse := fmt.Sprintf(`{"code":404,"mp3Url":""}`)
		w.Write([]byte(nullResponse))
		return
	}
	var trackDetail model.TrackDetail
	switch source {
	case "qq":
		trackDetail = qq.GetTrack(track)
	case "xiami":
		trackDetail = xiami.GetTrack(track)
	default:
		trackDetail = netease.GetTrack(track, country)
	}
	trackJson, _ := json.Marshal(trackDetail)
	w.Write([]byte(trackJson))
}
