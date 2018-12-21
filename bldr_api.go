package main

import (
	// "github.com/BurntSushi/toml"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type OriginKey struct {
	Origin   string `json:origin`
	Revision string `json:revision`
	Location string `json:location`
}

type BldrApi struct {
	url string
}

func (api BldrApi) fetchKeyPaths(origin string) []OriginKey {
	KEY_PATH := "/v1/depot/origins/" + origin + "/keys"
	log.Debug("Fetching all key paths for " + origin + " against bldr " + api.url)

	url := api.url + KEY_PATH
	log.Debug("HTTP GET " + url)

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var keys []OriginKey
	jsonErr := json.Unmarshal(body, &keys)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	log.Debug(keys)

	return keys
}

func (api BldrApi) fetchKeyData(key OriginKey) string {
	KEY_PATH := "/v1/depot" + key.Location
	log.Debug("Fetching key data for " + key.Origin + " rev " + key.Revision)

	url := api.url + KEY_PATH
	log.Debug("HTTP GET " + url)

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return string(body)
}
