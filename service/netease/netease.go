package netease

import (
	"MusicBox/model"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type MusicDetail struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

type netEaseSearch struct {
	Code   int             `json:"code"`
	Result ResultMusicType `json:result`
}

type ResultMusicType struct {
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

type TrackDetail struct {
	Code  int         `json:"code"`
	Track []TrackType `json:"data"`
}

type TrackType struct {
	Url string `json:"url"`
}

type netEasePlayList struct {
	Code   int                `json:"code"`
	Result ResultPlayListType `json:result`
}

type ResultPlayListType struct {
	PlayListCount int            `json:playlistCount`
	PlayLists     []PlayListType `json:playlists`
}
type PlayListType struct {
	CoverImgUrl string `json:"coverImgUrl"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
}

type netEasePlayListDetail struct {
	Code   int                      `json:"code"`
	Result ResultPlayListDetailType `json:result`
}

type ResultPlayListDetailType struct {
	TrackCount int        `json:trackCount`
	Tracks     []SongType `json:tracks`
}

//202.201.14.183
// http://music.163.com/api/playlist/detail?id=15451634
const (
	searchUrl   = "http://music.163.com/api/search/pc"
	playListUrl = "http://music.163.com/api/playlist/detail"
	trackUrl    = "http://music.163.com/weapi/song/enhance/player/url?csrf_token="
	cdnIP       = "202.201.14.183"
)

func SearchMusic(keyword string, limit int, page int) model.MusicSearch {
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
	var musicSearch model.MusicSearch
	musicSearch.Code = netEase.Code
	musicSearch.Data = musicDetail
	musicSearch.Total = netEase.Result.SongCount
	return musicSearch

}

func SearchPlayList(keyword string, limit int, page int) model.PlayListSearch {
	data := url.Values{}
	data.Set("s", keyword)
	data.Add("offset", strconv.Itoa(limit*(page-1)))
	data.Add("limit", strconv.Itoa(limit))
	data.Add("type", "1000")

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
	var netEase netEasePlayList
	json.Unmarshal(response, &netEase)

	playListDetail := make([]model.PlayListDetail, len(netEase.Result.PlayLists))
	for i := 0; i < len(playListDetail); i++ {
		playListDetail[i] = model.PlayListDetail{strconv.Itoa(netEase.Result.PlayLists[i].Id), netEase.Result.PlayLists[i].Name}
	}
	var playListSearch model.PlayListSearch
	playListSearch.Code = netEase.Code
	playListSearch.Data = playListDetail
	playListSearch.Total = netEase.Result.PlayListCount
	return playListSearch

}

func GetPlayList(id string) model.MusicSearch {
	data := url.Values{}
	data.Set("id", id)

	req, err := http.NewRequest("GET", playListUrl+"?"+data.Encode(), nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	var netEase netEasePlayListDetail
	json.Unmarshal(response, &netEase)

	musicDetail := make([]model.MusicDetail, netEase.Result.TrackCount)
	for i := 0; i < len(musicDetail); i++ {
		musicDetail[i] = model.MusicDetail{strconv.Itoa(netEase.Result.Tracks[i].Id), netEase.Result.Tracks[i].Name, netEase.Result.Tracks[i].Artist[0].Name, netEase.Result.Tracks[i].Album.Name}
	}
	var musicSearch model.MusicSearch
	musicSearch.Code = netEase.Code
	musicSearch.Data = musicDetail
	musicSearch.Total = netEase.Result.TrackCount
	return musicSearch

}

func GetTrack(id string, country string) model.TrackDetail {
	data := fmt.Sprintf(`{"ids":[%s],"br":320000,"csrf_token":""}`, id)
	req, err := http.NewRequest("POST", trackUrl, bytes.NewBufferString(encryptedRequest(data)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	var track TrackDetail
	json.Unmarshal(response, &track)
	var song model.TrackDetail
	if track.Code == 200 && track.Track[0].Url != "" {
		song = model.TrackDetail{track.Code, track.Track[0].Url}
	} else {
		song = model.TrackDetail{404, ""}
		return song
	}
	if country != "china" {
		song.Mp3Url = "http://" + cdnIP + "/" + strings.Split(song.Mp3Url, "http://")[1] + "?wshc_tag=0&&wsid_tag=b7fa6c76&wsiphost=ipdbm"
	}
	return song

}

func encryptedRequest(text string) string {
	var modulus = "00e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7b72" +
		"5152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280104e0312ecbd" +
		"a92557c93870114af6c9d05c4f7f0c3685b7a46bee255932575cce10b424d813cfe48" +
		"75d3e82047b97ddef52741d546b8e289dc6935b3ece0462db0a22b8e7"
	var nonce = "0CoJUm6Qyw8W8jud"
	var pubKey = "010001"
	var secKey = createSecretKey(16)
	encText := aesEncrypt([]byte(aesEncrypt([]byte(text), []byte(nonce))), []byte(secKey))
	encSecKey := rsaEncrypt(secKey, pubKey, modulus)
	data := url.Values{}
	data.Set("params", encText)
	data.Add("encSecKey", encSecKey)
	return data.Encode()

}

func createSecretKey(size int) string {
	choice := "012345679abcdef"
	result := ""
	for i := 0; i < size; i++ {
		result += string(choice[rand.Intn(len(choice))])
	}
	return result
}

// func aesEncrypt(text string, key string) string {
func aesEncrypt(text []byte, key []byte) string {
	iv := []byte("0102030405060708")
	plaintext := pad(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, len(plaintext))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[:], plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func rsaEncrypt(text string, key string, modulus string) string {
	text = reverse(text)
	text = hex.EncodeToString([]byte(text))
	base := new(big.Int)
	base.SetString(text, 16)
	exp := new(big.Int)
	exp.SetString(key, 16)
	mod := new(big.Int)
	mod.SetString(modulus, 16)
	result := new(big.Int)
	result.Exp(base, exp, nil)
	result.Mod(result, mod)
	return fmt.Sprintf("%0256x", result)
}
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
