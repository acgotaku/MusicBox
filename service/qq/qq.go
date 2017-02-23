package qq

// http://i.y.qq.com/s.music/fcgi-bin/search_for_qq_cp?format=json&platform=h5&w=tsubasa&n=20&p=1
import (
	"MusicBox/service"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type QQMusicSearch struct {
	Code int      `json:"code"`
	Data DataType `json:data`
}

type DataType struct {
	Keyword string   `json:"keyword"`
	Song    SongType `json:song`
}
type SongType struct {
	Totalnum int            `json:"totalnum"`
	List     []SongListType `json:"list"`
}
type SongListType struct {
	Id     string       `json:"songmid"`
	Name   string       `json:"songname"`
	Mp3Url string       `json:"mp3Url"`
	Artist []ArtistType `json:"singer"`
	Album  string       `json:"albumname"`
}

type ArtistType struct {
	Name string `json:"name"`
}

const searchUrl = "http://i.y.qq.com/s.music/fcgi-bin/search_for_qq_cp"

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data := url.Values{}
	data.Set("w", r.URL.Query().Get("keyword"))
	data.Add("format", "json")
	data.Add("platform", "h5")
	data.Add("n", "20")
	data.Add("p", "1")
	req, err := http.NewRequest("GET", searchUrl+"?"+data.Encode(), nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	var qqMusic QQMusicSearch
	json.Unmarshal(response, &qqMusic)
	musicDetail := make([]service.MusicDetail, len(qqMusic.Data.Song.List))
	for i := 0; i < len(musicDetail); i++ {
		musicDetail[i] = service.MusicDetail{qqMusic.Data.Song.List[i].Id, qqMusic.Data.Song.List[i].Name, qqMusic.Data.Song.List[i].Artist[0].Name, qqMusic.Data.Song.List[i].Album}
	}
	music, _ := json.Marshal(musicDetail)
	w.Write(music)
}

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data := url.Values{}
	data.Set("w", r.URL.Query().Get("keyword"))
	data.Add("format", "json")
	data.Add("platform", "h5")
	data.Add("n", "20")
	data.Add("p", "1")
	req, err := http.NewRequest("GET", searchUrl, bytes.NewBufferString(data.Encode()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)

	w.Write(response)

}
