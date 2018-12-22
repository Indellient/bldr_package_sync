package main

import (
	"strings"
	"testing"
)

func TestFetchKeyPaths(t *testing.T) {
	api := BldrApi{Url: "https://bldr.habitat.sh"}
	if len(api.fetchKeyPaths("core")) <= 0 {
		t.Error("Fetching Paths returned an slice <= 0")
	}
}

func TestFetchKeyData(t *testing.T) {
	api := BldrApi{Url: "https://bldr.habitat.sh"}
	key := api.fetchKeyPaths("core")[0]
	data := api.fetchKeyData(key)

	// Check that the returned key is not empty
	if data == "" {
		t.Error("Fetching Key data failed")
	}

	// Test that the returned key contains the origin and revision
	if !strings.Contains(data, key.Origin+"-"+key.Revision) {
		t.Error("Fetching Key data failed, contains likely not valid")
	}
}
