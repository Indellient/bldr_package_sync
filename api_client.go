package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Implement status codes: https://golang.org/src/net/http/status.go
func performGetRequest(url string) (*http.Response, error) {

	sleepTime := 0
	retryCount := 0

	res, err := getRequest(url)

	for err != nil {
		if retryCount >= 10 {
			log.Error("Request was tried more than 10 times, giving up: " + url)
			return nil, err
		}
		sleepTime = sleepTime + 1
		retryCount = retryCount + 1
		time.Sleep(time.Duration(sleepTime) * time.Second)
		res, err = getRequest(url)
	}

	return res, nil
}

// Performs an HTTP GET Method
func getRequest(url string) (*http.Response, error) {

	client := http.Client{
		Timeout: time.Second * 300,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithError(err).Error("Error in GET request")
		return nil, err
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.WithError(err).Error("Error in GET request")
		return nil, err
	}

	if err != nil {
		log.WithError(err).Error("failed to listen for messages")
		return nil, err
	}

	if res.StatusCode == http.StatusTooManyRequests {
		err = errors.New("Too many requests being made against host, reporting error and trying again")
		log.WithError(err).Error("Failed to perform GET")
		return nil, err
	}

	if res.StatusCode > http.StatusMultipleChoices {
		err = fmt.Errorf("Incorrect response code returned, received: %v", res.StatusCode)
		log.WithError(err).Error("Failed to perform GET")
		return nil, err
	}

	return res, nil
}
