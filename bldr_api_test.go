package main

import (
	"strings"
	"testing"
)

const INDELLIENT_PKG_COUNT = 18

func TestListPackagesAndDeps(t *testing.T) {
	api := BldrApi{Url: "https://bldr.habitat.sh"}
	pkgs := api.listAllPackages("indellient", "stable")

	if len(pkgs.Data) <= 0 {
		t.Error("Fetching Packages returned an slice <= 0")
	}

	if len(pkgs.Data) < INDELLIENT_PKG_COUNT {
		t.Error("Fetching All Packages failed, highly doubt packages from upstream were deleted ", len(pkgs.Data))
	}

	for _, p := range pkgs.Data {
		deps := api.fetchPackageDeps(p)
		if len(deps) != 0 {
			for _, d := range deps {
				if d.Name == "" {
					t.Error("Fetching Packages returned an slice <= 0", p)
				}
			}
		}
	}
}

func TestListPackages(t *testing.T) {
	api := BldrApi{Url: "https://bldr.habitat.sh"}
	pkgs := api.listAllPackages("indellient", "stable")

	if len(pkgs.Data) <= 0 {
		t.Error("Fetching Packages returned an slice <= 0")
	}

	if len(pkgs.Data) < INDELLIENT_PKG_COUNT {
		t.Error("Fetching All Packages failed, highly doubt packages from upstream were deleted ", len(pkgs.Data))
	}
}

func TestFetchPackageDeps(t *testing.T) {
	api := BldrApi{Url: "https://bldr.habitat.sh"}
	data := PackageData{Origin: "core", Name: "airlock", Version: "10", Release: "20171027222310"}
	pkg := api.fetchPackageDeps(data)

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
