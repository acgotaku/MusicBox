package qq

// http://i.y.qq.com/s.music/fcgi-bin/search_for_qq_cp?format=json&platform=h5&w=tsubasa&n=20&p=1
import (
	"MusicBox/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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

type TrackDetail struct {
	Code int    `json:"code"`
	Key  string `json:"key"`
}

const searchUrl = "http://i.y.qq.com/s.music/fcgi-bin/search_for_qq_cp"

const trackUrl = "http://base.music.qq.com/fcgi-bin/fcg_musicexpress.fcg?json=3&format=json&guid=780782017"

func SearchMusic(keyword string, limit int, page int) []model.MusicDetail {
	data := url.Values{}
	data.Set("w", keyword)
	data.Add("format", "json")
	data.Add("platform", "h5")
	data.Add("n", strconv.Itoa(limit))
	data.Add("p", strconv.Itoa(page))
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
	musicDetail := make([]model.MusicDetail, len(qqMusic.Data.Song.List))
	for i := 0; i < len(musicDetail); i++ {
		musicDetail[i] = model.MusicDetail{qqMusic.Data.Song.List[i].Id, qqMusic.Data.Song.List[i].Name, qqMusic.Data.Song.List[i].Artist[0].Name, qqMusic.Data.Song.List[i].Album}
	}
	return musicDetail
}

// http://base.music.qq.com/fcgi-bin/fcg_musicexpress.fcg?json=3&format=json
func GetTrack(id string) model.TrackDetail {
	req, err := http.NewRequest("GET", trackUrl, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	var track TrackDetail
	json.Unmarshal(response, &track)
	//var url = "http://cc.stream.qqmusic.qq.com/C200" +  track.id.slice('qqtrack_'.length)  + ".m4a?vkey=" +token + "&fromtag=0&guid=780782017";
	url := fmt.Sprintf(`http://cc.stream.qqmusic.qq.com/C200%s.m4a?vkey=%s&fromtag=0&guid=780782017`, id, track.Key)
	trackDetail := model.TrackDetail{200, url}
	return trackDetail
}
