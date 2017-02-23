package xiami

import (
	"MusicBox/model"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const searchUrl = "http://music.163.com/api/search/pc"

type netEaseSearch struct {
	Code   int        `json:"code"`
	Result ResultType `json:result`
}

type ResultType struct {
	QueryCorrected []string   `json:"queryCorrected"`
	SongCount      int        `json:songCount`
	Songs          []SongType `json:songs`
}
type SongType struct {
	Id     int          `json:"id"`
	Name   string       `json:"name"`
	Mp3Url string       `json:"mp3Url"`
	Artist []ArtistType `json:"artists"`
	Album  AlbumType    `json:"album"`
}

type AlbumType struct {
	Name   string `json:"name"`
	ImgUrl string `json:"picUrl"`
}
type ArtistType struct {
	Name string `json:"name"`
}

func SearchMusic(keyword string, limit int, page int) []model.MusicDetail {
	data := url.Values{}
	data.Set("s", keyword)
	data.Add("offset", strconv.Itoa(limit*(page-1)))
	data.Add("limit", strconv.Itoa(limit))
	data.Add("type", "1")
	req, err := http.NewRequest("POST", searchUrl, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	var netEase netEaseSearch
	json.Unmarshal(response, &netEase)
	musicDetail := make([]model.MusicDetail, len(netEase.Result.Songs))
	for i := 0; i < len(musicDetail); i++ {
		musicDetail[i] = model.MusicDetail{strconv.Itoa(netEase.Result.Songs[i].Id), netEase.Result.Songs[i].Name, netEase.Result.Songs[i].Artist[0].Name, netEase.Result.Songs[i].Album.Name}
	}
	return musicDetail

}
