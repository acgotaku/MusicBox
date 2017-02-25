package model

type MusicSearch struct {
	Code  int           `json:"code"`
	Data  []MusicDetail `json:"data"`
	Total int           `json:"total"`
}

type MusicDetail struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}
type TrackDetail struct {
	Code   int    `json:"code"`
	Mp3Url string `json:"mp3Url"`
}
