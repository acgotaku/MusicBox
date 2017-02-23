package netease

import (
	"MusicBox/service"
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
)

type MusicDetail struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

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

const searchUrl = "http://music.163.com/api/search/pc"

const trackUrl = "http://music.163.com/weapi/song/enhance/player/url?csrf_token="

func SearchMusic(keyword string, limit int, page int) []service.MusicDetail {
	data := url.Values{}
	data.Set("s", keyword)
	data.Add("offset", limit*(page-1))
	data.Add("limit", limit)
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
	musicDetail := make([]service.MusicDetail, len(netEase.Result.Songs))
	for i := 0; i < len(musicDetail); i++ {
		musicDetail[i] = service.MusicDetail{strconv.Itoa(netEase.Result.Songs[i].Id), netEase.Result.Songs[i].Name, netEase.Result.Songs[i].Artist[0].Name, netEase.Result.Songs[i].Album.Name}
	}
	return musicDetail

}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data := url.Values{}
	data.Set("s", r.URL.Query().Get("keyword"))
	data.Add("offset", "10")
	data.Add("limit", "20")
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
	musicDetail := make([]service.MusicDetail, len(netEase.Result.Songs))
	for i := 0; i < len(musicDetail); i++ {
		musicDetail[i] = service.MusicDetail{strconv.Itoa(netEase.Result.Songs[i].Id), netEase.Result.Songs[i].Name, netEase.Result.Songs[i].Artist[0].Name, netEase.Result.Songs[i].Album.Name}
	}
	music, _ := json.Marshal(musicDetail)
	w.Write(music)

}

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data := fmt.Sprintf(`{"ids":[%s],"br":320000,"csrf_token":""}`, r.URL.Query().Get("id"))
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
	// var netEase netEaseSearch
	// json.Unmarshal(response, &netEase)
	w.Write(response)
	// fmt.Printf("%+v\n", netEase)
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