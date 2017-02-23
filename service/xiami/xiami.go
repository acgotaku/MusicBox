package xiami

import (
	"MusicBox/model"
	"encoding/json"

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
