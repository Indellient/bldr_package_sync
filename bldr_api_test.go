package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestListPackagesAndDeps(t *testing.T) {
	api := BldrApi{Url: "https://bldr.habitat.sh"}
	pkgs := api.listAllPackages("skylerto", "stable")

	if len(pkgs.Data) <= 0 {
		t.Error("Fetching Packages returned an slice <= 0")
	}

	if len(pkgs.Data) < 4737 {
		t.Error("Fetching All Packages failed, highly doubt packages from core were deleted")
	}

	for _, p := range pkgs.Data {
		deps := api.fetchPackageDeps(p)
		if len(deps) <= 0 {
			t.Error("Fetching Packages returned an slice <= 0")
		}
	}
}

func TestListPackages(t *testing.T) {
	api := BldrApi{Url: "https://bldr.habitat.sh"}
	pkgs := api.listAllPackages("core", "stable")

	if len(pkgs.Data) <= 0 {
		t.Error("Fetching Packages returned an slice <= 0")
	}

	if len(pkgs.Data) < 4737 {
		t.Error("Fetching All Packages failed, highly doubt packages from core were deleted")
	}
}

func TestFetchPackageDeps(t *testing.T) {
	api := BldrApi{Url: "https://bldr.habitat.sh"}
	data := PackageData{Origin: "core", Name: "airlock", Version: "10", Release: "20171027222310"}
	pkg := api.fetchPackageDeps(data)
	fmt.Println(pkg)

	if len(pkg) <= 0 {
		t.Error("Error fetching package info")
	}
}

func TestFetchPackageData(t *testing.T) {
	api := BldrApi{Url: "https://bldr.habitat.sh"}
	data := PackageData{Origin: "core", Name: "airlock", Version: "10", Release: "20171027222310"}
	pkg := api.fetchPackage(data)

	if pkg.Name != "airlock" {
		t.Error("Error fetching package info")
	}

}

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
