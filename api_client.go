package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var SLEEP_TIME = 0

// Implement status codes: https://golang.org/src/net/http/status.go
func performGetRequest(url string) *http.Response {

	res := getRequest(url)

	if res.StatusCode == http.StatusTooManyRequests {
		log.Error("Too many requests being made against host, increasing the sleep time and performaing the request again", res.StatusCode)
		SLEEP_TIME = SLEEP_TIME + 2
		return getRequest(url)
	}

	if res.StatusCode > http.StatusMultipleChoices {
		log.Error("Incorrect response code returned ", res.StatusCode)
		return nil
	}

	return res
}

// Performs an HTTP GET Method
func getRequest(url string) *http.Response {

	log.Debug("HTTP GET " + url)

	time.Sleep(time.Second * time.Duration(SLEEP_TIME))

	client := http.Client{
		Timeout: time.Second * 300,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	return res
}
