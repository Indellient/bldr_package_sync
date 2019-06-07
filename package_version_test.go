package main

import (
	"testing"
)

func TestVersionsNoPatch(t *testing.T) {
	packageData := PackageData{
		Origin:  "indellient",
		Name:    "consul",
		Version: "1.4.3",
		Release: "20190327191349",
	}

	v := "<= 1.4"
	if packageData.MatchesVersion(v) {
		t.Error("Package Version test Fails to validate contraint " + v)
	}

	v = ">= 1.4"
	if !packageData.MatchesVersion(v) {
		t.Error("Package Version test Fails to validate contraint " + v)
	}

	v = "= 1.4"
	if packageData.MatchesVersion(v) {
		t.Error("Package Version test Fails to validate contraint " + v)
	}

	v = "< 1.4"
	if packageData.MatchesVersion(v) {
		t.Error("Package Version test Fails to validate contraint " + v)
	}

	v = "> 1.4"
	if !packageData.MatchesVersion(v) {
		t.Error("Package Version test Fails to validate contraint " + v)
	}

	v = "~> 1.4"
	if !packageData.MatchesVersion(v) {
		t.Error("Package Version test Fails to validate contraint " + v)
	}
}

func TestNewPackageDataFullString(t *testing.T) {
	packageData := PackageData{
		Origin:  "indellient",
		Name:    "consul",
		Version: "1.4.3",
		Release: "20190327191349",
	}
	ident := NewPackageData("indellient/consul/1.4.3/20190327191349")

	if ident != packageData {
		t.Error("Failed to parse new packge ident")
	}
}

func TestNewPackageDataOriginName(t *testing.T) {
	packageData := PackageData{
		Origin: "indellient",
		Name:   "consul",
	}
	ident := NewPackageData("indellient/consul")

	if ident != packageData {
		t.Error("Failed to parse new packge ident")
	}
}

func TestPackageConstraintValidatePackageContraint(t *testing.T) {
	packageData := PackageData{
		Origin:  "indellient",
		Name:    "consul",
		Version: "1.4.3",
		Release: "20190327191349",
	}

	packageConstraint := PackageConstraint{
		Name:       "indellient/consul",
		Constraint: "> 1.4",
	}

	if !packageConstraint.ValidatePackageContraint(packageData) {
		t.Errorf("Failed to validate contraint properly %v, %v", packageData, packageConstraint)
	}

	packageConstraint = PackageConstraint{
		Name:       "indellient/consul",
		Constraint: "> 2.0.0",
	}

	if packageConstraint.ValidatePackageContraint(packageData) {
		t.Errorf("Failed to validate contraint properly %v, %v", packageData, packageConstraint)
	}
}
