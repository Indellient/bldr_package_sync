package main

import (
	// "github.com/BurntSushi/toml"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"
)

type OriginKey struct {
	Origin   string `json:origin`
	Revision string `json:revision`
	Location string `json:location`
}

type BldrApi struct {
	Url       string `toml:url`
	AuthToken string `toml:authToken`
}

type PackageData struct {
	Origin  string `json:origin`
	Name    string `json:name`
	Version string `json:version`
	Release string `json:release`
}

type Packages struct {
	Start int           `json:"range_start"`
	End   int           `json:"range_end"`
	Total int           `json:"total_count"`
	Data  []PackageData `json:data`
}

type Package struct {
	Channels   []string      `json:"channels"`
	Checksum   string        `json:"checksum"`
	Config     string        `json:"config"`
	CreatedAt  string        `json:"created_at"`
	Deps       []PackageData `json:"deps"`
	TDeps      []PackageData `json:"tdeps"`
	Exposes    []int         `json:"exposes"`
	Id         string        `json:"id"`
	Ident      PackageData   `json:"ident"`
	IdentArray []string      `json:"ident_array"`
	IsAService bool          `json:"is_a_service"`
	Manifest   string        `json:"manifest"`
	Name       string        `json:"name"`
	Origin     string        `json:"origin"`
	OwnerId    string        `json:"owner_id"`
	Target     string        `json:"target"`
	UpdatedAt  string        `json:"updated_at"`
	Visibility string        `json:"visibility"`
}

// Package dependencies are allows in the stable channel
// Therefore we should never include the package we're dealing with
// in its tdeps array
func (api BldrApi) fetchPackageDeps(pkg PackageData) []PackageData {
	var pkgs []PackageData

	data := api.fetchPackage(pkg)

	if len(data.TDeps) <= 0 {
		return pkgs
	}

	for _, p := range data.TDeps {
		if len(pkgs) <= 0 {
			pkgs = []PackageData{p}
		}
		return append(pkgs, api.fetchPackageDeps(p)...)
	}

	return pkgs
}

func (api BldrApi) fetchPackage(pkg PackageData) Package {
	var data Package
	pkgName := fmt.Sprintf("%s/%s/%s/%s", pkg.Origin, pkg.Name, pkg.Version, pkg.Release)

	url := fmt.Sprintf("%s/v1/depot/pkgs/%s", api.Url, pkgName)

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

	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return data
}

func (api BldrApi) listAllPackages(origin string, channel string) Packages {
	packages := api.listPackages(origin, channel)
	resultsPerPage := float64(packages.End - packages.Start)
	iterations := math.Ceil(float64(packages.Total) / resultsPerPage)
	for i := float64(1); i <= iterations; i++ {
		pkgs := api.listPackagesRange(origin, channel, int(i*resultsPerPage))
		packages.Data = append(packages.Data, pkgs.Data...)
	}
	return packages
}

func (api BldrApi) listPackages(origin string, channel string) Packages {
	PACKGE_PATH := "/v1/depot/channels/" + origin + "/" + channel + "/pkgs"

	url := api.Url + PACKGE_PATH
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

	var pkgs Packages
	jsonErr := json.Unmarshal(body, &pkgs)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return pkgs
}

func (api BldrApi) listPackagesRange(origin string, channel string, count int) Packages {
	// PACKGE_PATH := "/v1/depot/channels/" + origin + "/" + channel + "/pkgs?range=" + count
	PACKGE_PATH := fmt.Sprintf("/v1/depot/channels/%s/%s/pkgs?range=%d", origin, channel, count)

	url := api.Url + PACKGE_PATH
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

	var pkgs Packages
	jsonErr := json.Unmarshal(body, &pkgs)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return pkgs
}

func (api BldrApi) fetchKeyPaths(origin string) []OriginKey {
	KEY_PATH := "/v1/depot/origins/" + origin + "/keys"

	url := api.Url + KEY_PATH
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

	return keys
}

func (api BldrApi) fetchKeyData(key OriginKey) string {
	KEY_PATH := "/v1/depot" + key.Location

	url := api.Url + KEY_PATH
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

func (api BldrApi) uploadOriginKey(filename string, key string, origin string) bool {

	dir := os.TempDir()
	file := dir + filename

	if err := ioutil.WriteFile(file, []byte(key), 0777); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	log.Debug("Created File: " + file)

	importPublicKey(api, dir, file)

	os.Remove(file)
	return true
}

func difference(upstream []OriginKey, target []OriginKey) []OriginKey {
	var diff []OriginKey

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	// for i := 0; i < 1; i++ {
	for _, s1 := range upstream {
		found := false
		for _, s2 := range target {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, s1)
		}
		// }
		// Swap the slices, only if it was the first loop
		// if i == 0 {
		// 	slice1, slice2 = slice2, slice1
		// }
	}

	return diff
}
