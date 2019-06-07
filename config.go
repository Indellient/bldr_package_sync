package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// FeatureList is a list of all the available features
var FeatureList = [...]string{"PACKAGE_CONSTRAINTS"}

// Origin is a Definition of a habitat origin as interpreted from
// configuration files.
type Origin struct {
	Name     string   `toml:"name"`
	Channels []string `toml:"channels"`
}

// PackageConstraint defines a contrainst of a package and it's versions to upload
type PackageConstraint struct {
	Name       string `toml:"name"`
	Constraint string `toml:"constraint"`
}

// Structure of a configuration file.
type Config struct {
	Upstream           BldrApi
	LogLevel           string `toml:"log_level"`
	TempDir            string `toml:"temp_dir"`
	Target             BldrApi
	Env                []string
	Features           []string
	Interval           int
	Origins            []Origin            `toml:"origin"`
	PackageConstraints []PackageConstraint `toml:"package"`
}

func (config Config) ValidatePackageContraint(packageData PackageData) bool {
	for _, pkg := range config.PackageConstraints {
		if !pkg.ValidatePackageContraint(packageData) {
			return false
		}
	}
	return true
}

func (packageConstraint PackageConstraint) ValidatePackageContraint(packageData PackageData) bool {
	// If we've made it through checking all constraints without changing we should upload
	result := true
	if fmt.Sprintf("%s/%s", packageData.Origin, packageData.Name) == packageConstraint.Name {
		if !packageData.MatchesVersion(packageConstraint.Constraint) {
			result = false
		}
	}

	log.WithFields(log.Fields{
		"packageConstraint": packageConstraint,
		"packageData":       packageData,
		"result":            result,
	}).Debugf("Validating Package Constraint")
	return result
}

func (config Config) PackageContraintEnabled() bool {
	return contains(config.Features, "PACKAGE_CONSTRAINTS")
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
