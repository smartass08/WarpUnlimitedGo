package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	referrer = ""

	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	allnumbers = []rune("01234456789")

	warpHeaders = http.Header{
		"Content-Type": []string{"application/json; charset=UTF-8"},
		"Host": []string{"api.cloudflareclient.com"},
		"Connection": []string{"Keep-Alive"},
		"Accept-Encoding": []string{"gzip"},
		"User-Agent": []string{"okhttp/3.12.1"},
	}
	DefaultTransport = &http.Transport{
		// Match app's TLS config or API will reject us with code 403 error 1020
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS12},
		ForceAttemptHTTP2: false,
		// From http.DefaultTransport
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
)


type WarpBody struct {
	Key         string `json:"key"`
	InstallId   string `json:"install_id"`
	FcmToken    string `json:"fcm_token"`
	Referrer    string `json:"referrer"`
	WarpEnabled bool   `json:"warp_enabled"`
	Tos         string `json:"tos"`
	Type        string `json:"type"`
	Locale		string  `default:"es_ES" json:"locale"`
}

func randSeq(n int, issnumber bool) string {
	b := make([]rune, n)
	if issnumber{
		for i := range b {
			b[i] = allnumbers[rand.Intn(len(allnumbers))]
		}
		return string(b)
	}
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateString(length int, onlyNumbers bool) string {
	rand.Seed(time.Now().UnixNano())
	return randSeq(length, onlyNumbers)
}

func GetTimestamp() string {
	return getTimestamp(time.Now())
}

func getTimestamp(t time.Time) string {
	timestamp := t.Format(time.RFC3339Nano)
	return timestamp
}

func doRequest() (*http.Response, error) {
	tempWarpID := generateString(22, false)
	cfurI := fmt.Sprintf("https://api.cloudflareclient.com/v0a%v/reg", generateString(3, true))
	warpBody := WarpBody{
		Key:         fmt.Sprintf("%v=", generateString(43, false)),
		InstallId:   tempWarpID,
		FcmToken:    fmt.Sprintf("%v:APA91b%v", tempWarpID, generateString(134, false)),
		Referrer:    referrer,
		WarpEnabled: false,
		Tos:       	GetTimestamp(),
		Type:        "Android",
		Locale:      "es_ES",
	}
	jsondata, err:= json.Marshal(warpBody)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	jsonRaw := bytes.NewReader(jsondata)
	r, err := http.NewRequest(http.MethodPost,cfurI, jsonRaw)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	r.Header = warpHeaders
	client := &http.Client{Transport: DefaultTransport}
	res, err := client.Do(r)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	return res, nil
}

func main()  {
	count := 0
	for{
		info, err := doRequest()
		if err != nil {
			log.Println("Error occurred while making the request")
			continue
		}
		if info.StatusCode == 200{
			count += 1
			log.Printf("Success! 1 GB added, Total added := %v GB\n", count)
		} else {
			log.Printf("Error occurred while making the request, Response Code := %v\nretrying...\n", info.StatusCode)
		}

		time.Sleep(time.Second*7)
	}
}
