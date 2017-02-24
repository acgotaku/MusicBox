package xiami

import (
	"MusicBox/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// http://api.xiami.com/web?v=2.0&app_key=1&key=tsubasa&page=1&limit=20&r=search/songs
const searchUrl = "http://api.xiami.com/web"

type XiamiSearch struct {
	Code int      `json:"state"`
	Data DataType `json:data`
}

type DataType struct {
	Totalnum int            `json:"total"`
	Songs    []SongListType `json:songs`
}
type SongListType struct {
	Id     int    `json:"song_id"`
	Name   string `json:"song_name"`
	Mp3Url string `json:"listen_file"`
	Artist string `json:"artist_name"`
	Album  string `json:"album_name"`
}

type TrackDetail struct {
	Code  int       `json:"state"`
	Track TrackType `json:"data"`
}
type TrackType struct {
	Song SongType `json:"song"`
}

type SongType struct {
	Url string `json:"listen_file"`
}

func SearchMusic(keyword string, limit int, page int) []model.MusicDetail {
	data := url.Values{}
	data.Set("key", keyword)
	data.Add("app_key", "1")
	data.Add("page", strconv.Itoa(page))
	data.Add("limit", strconv.Itoa(limit))
	data.Add("r", "search/songs")
	req, err := http.NewRequest("GET", searchUrl+"?"+data.Encode(), nil)
	req.Header.Set("Referer", "http://m.xiami.com/")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	var xiamiMusic XiamiSearch
	json.Unmarshal(response, &xiamiMusic)
	musicDetail := make([]model.MusicDetail, len(xiamiMusic.Data.Songs))
	for i := 0; i < len(musicDetail); i++ {
		musicDetail[i] = model.MusicDetail{strconv.Itoa(xiamiMusic.Data.Songs[i].Id), xiamiMusic.Data.Songs[i].Name, xiamiMusic.Data.Songs[i].Artist, xiamiMusic.Data.Songs[i].Album}
	}
	return musicDetail

}

// http://api.xiami.com/web?v=2.0&app_key=1&r=song/detail&id=1769150238
func GetTrack(id string) model.TrackDetail {
	trackUrl := fmt.Sprintf(`http://api.xiami.com/web?v=2.0&app_key=1&r=song/detail&id=%s`, id)
	req, err := http.NewRequest("GET", trackUrl, nil)
	req.Header.Set("Referer", "http://m.xiami.com/")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	var track TrackDetail
	json.Unmarshal(response, &track)

	trackDetail := model.TrackDetail{200, track.Track.Song.Url}
	return trackDetail
}
