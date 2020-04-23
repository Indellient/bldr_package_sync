package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	version "github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
)

const BAD_CODE = 300

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

func (api BldrApi) downloadPackage(pack Package) string {

	pkg := pack.Ident
	pkgName := fmt.Sprintf("%s/%s/%s/%s", pkg.Origin, pkg.Name, pkg.Version, pkg.Release)
	url := fmt.Sprintf("%s/v1/depot/pkgs/%s/download?target=%s", api.Url, pkgName, pack.Target)

	dir := config.TempDir
	if config.TempDir == "" {
		dir = os.TempDir()
	}

	hartFile := fmt.Sprintf("%s-%s-%s-%s-%s.hart", pkg.Origin, pkg.Name, pkg.Version, pkg.Release, pack.Target)
	location := filepath.Join(dir, hartFile)
	log.Debug("Downloading to file ", location)

	// Create the file
	out, err := os.Create(location)
	if err != nil {
		log.Error(err)
	}
	defer out.Close()

	resp := performGetRequest(url)
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error(err)
	}

	return location
}

// Package dependencies are allows in the stable channel
// Therefore we should never include the package we're dealing with
// in its tdeps array
func (api BldrApi) fetchPackageDeps(pkg PackageData) ([]PackageData, error) {
	data, err := api.fetchPackage(pkg)
	if err != nil {
		return []PackageData{}, err
	}
	return data.TDeps, nil
}

func (api BldrApi) packageExists(pkg PackageData) bool {
	pkgName := fmt.Sprintf("%s/%s/%s/%s", pkg.Origin, pkg.Name, pkg.Version, pkg.Release)

	url := fmt.Sprintf("%s/v1/depot/pkgs/%s", api.Url, pkgName)

	client := http.Client{
		Timeout: time.Second * 30, // Maximum of 30 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	return res.StatusCode == http.StatusOK

}

func (api BldrApi) fetchPackage(pkg PackageData) (Package, error) {
	var data Package
	pkgName := fmt.Sprintf("%s/%s/%s/%s", pkg.Origin, pkg.Name, pkg.Version, pkg.Release)

	url := fmt.Sprintf("%s/v1/depot/pkgs/%s", api.Url, pkgName)

	res := performGetRequest(url)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Error(readErr)
		return Package{}, readErr
	}

	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Error(jsonErr)
		return Package{}, jsonErr
	}

	return data, nil
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
	res := performGetRequest(url)

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
	PACKGE_PATH := fmt.Sprintf("/v1/depot/channels/%s/%s/pkgs?range=%d", origin, channel, count)

	url := api.Url + PACKGE_PATH
	res := performGetRequest(url)

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
	res := performGetRequest(url)

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
	res := performGetRequest(api.Url + "/v1/depot" + key.Location)
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return string(body)
}

func (api BldrApi) uploadOriginKey(filename string, key string, origin string) bool {

	dir := config.TempDir
	if config.TempDir == "" {
		dir = os.TempDir()
	}
	file := filepath.Join(dir, filename)

	if err := ioutil.WriteFile(file, []byte(key), 0777); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	log.Debug("Created File: " + file)

	importPublicKey(api, dir, file)
	return true
}

func packageDifference(upstream []PackageData, target []PackageData) []PackageData {
	var diff []PackageData

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
	}

	return diff
}

func difference(upstream []OriginKey, target []OriginKey) []OriginKey {
	var diff []OriginKey

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
	}

	return diff
}

func NewPackageData(data string) PackageData {
	var packageData PackageData
	identArray := strings.Split(data, "/")
	for i := 0; i < len(identArray); i++ {
		if i == 0 {
			packageData.Origin = identArray[0]
		}

		if i == 1 {
			packageData.Name = identArray[1]
		}

		if i == 2 {
			packageData.Version = identArray[2]
		}

		if i == 3 {
			packageData.Release = identArray[3]
		}
	}
	return packageData
}

// MatchesVersion matches based on the following contraints
// "":   constraintEqual,
// "=":  constraintEqual,
// "!=": constraintNotEqual,
// ">":  constraintGreaterThan,
// "<":  constraintLessThan,
// ">=": constraintGreaterThanEqual,
// "<=": constraintLessThanEqual,
// "~>": constraintPessimistic,
func (packageData PackageData) MatchesVersion(constraint string) bool {
	v1, err := version.NewVersion(packageData.Version)
	if err != nil {
		log.Error(err)
		return false
	}

	constraints, err := version.NewConstraint(constraint)
	if err != nil {
		log.Error(err)
		return false
	}

	log.WithFields(log.Fields{
		"version":    v1,
		"constraint": constraint,
		"statisfies": constraints.Check(v1),
	}).Debugf("Checking Version Contraint %s %s", v1, constraint)
	return constraints.Check(v1)
}
