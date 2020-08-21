package main

import (
	// "github.com/BurntSushi/toml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	log "github.com/sirupsen/logrus"
)

// Syncer holds the configuration values for the sync
type Syncer struct {
	config Config
}

func (syncer Syncer) syncPackages(origin string, channel string, upstream BldrApi, target BldrApi) bool {
	log.Info("Beginning the package sync process")

	upstreamPkgsChan := make(chan Packages)
	go func() {
		pkgs := upstream.listAllPackages(origin, channel)
		log.Info(fmt.Sprintf("Found %d packages on %s", len(pkgs.Data), upstream.Url))
		upstreamPkgsChan <- pkgs
	}()

	targetPkgsChan := make(chan Packages)
	go func() {
		targetPkgs := target.listAllPackages(origin, channel)
		log.Info(fmt.Sprintf("Found %d packages on %s", len(targetPkgs.Data), target.Url))
		targetPkgsChan <- targetPkgs
	}()

	upstreamPkgs := <-upstreamPkgsChan
	targetPkgs := <-targetPkgsChan

	// Good enough to figure out the difference before calculating deps
	pkgDatas := packageDifference(upstreamPkgs.Data, targetPkgs.Data)

	log.Info(fmt.Sprintf("Determining TDEPS for %d packages", len(pkgDatas)))

	// Currently adding multi-thread support for syncing packages pounds both upstream and target
	// APIs, typically resulting in Fatal API calls.
	for j, p := range pkgDatas {

		if syncer.config.PackageContraintEnabled() {
			log.Info("Validating package against contraints")
			if syncer.config.ValidatePackageContraint(p) {
				log.Info(fmt.Sprintf("package [%d/%d]", j, len(pkgDatas)))
				syncer.syncPackage(upstream, target, p, channel)
			}
		} else {
			log.Info(fmt.Sprintf("package [%d/%d]", j, len(pkgDatas)))
			syncer.syncPackage(upstream, target, p, channel)
		}

	}

	return true
}

func (syncer Syncer) syncKeys(origin string, upstream BldrApi, target BldrApi) bool {
	log.Info("Beginning the key sync process")
	upstreamKeys := upstream.fetchKeyPaths(origin)
	targetKeys := target.fetchKeyPaths(origin)

	for _, upstreamKey := range upstreamKeys {
		keyData := upstream.fetchKeyData(upstreamKey)
		keyFileName := upstreamKey.Origin + "-" + upstreamKey.Revision + ".pub"
		if err := ioutil.WriteFile(path.Join(syncer.config.TempDir, keyFileName), []byte(keyData), 0777); err != nil {
			log.Fatal("Failed to write to temporary file", err)
		}

	}

	keys := difference(upstreamKeys, targetKeys)
	log.Debug("Uploading diffed keys")
	log.Debug(keys)

	for _, key := range keys {
		data := upstream.fetchKeyData(key)
		log.Debug(data)
		fileName := key.Origin + "-" + key.Revision + ".pub"
		target.uploadOriginKey(fileName, data, key.Origin)
	}

	return true
}

func (syncer Syncer) run() error {
	for {
		for _, origin := range syncer.config.Origins {
			syncer.syncKeys(origin.Name, syncer.config.Upstream, syncer.config.Target)
		}

		for _, origin := range syncer.config.Origins {
			for _, channel := range origin.Channels {
				log.Info(fmt.Sprintf("Syncing packages for %s on channel %s", origin.Name, channel))
				syncer.syncPackages(origin.Name, channel, syncer.config.Upstream, syncer.config.Target)
			}
		}
		log.Info(fmt.Sprintf("Sync process finished, Sleeping for %d seconds", config.Interval))
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

func (syncer Syncer) syncPackage(upstream BldrApi, target BldrApi, p PackageData, channel string) {
	files := []string{}
	log.Info(fmt.Sprintf("Syncing package %v", p))
	deps, err := upstream.fetchPackageDeps(p)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(fmt.Sprintf("Determined deps %s", deps))
	for i, pkg := range deps {
		pkgName := fmt.Sprintf("%s/%s/%s/%s", pkg.Origin, pkg.Name, pkg.Version, pkg.Release)

		log.Info(fmt.Sprintf("Dependancy [%d/%d] %s", i+1, len(deps), pkgName))
		if !target.packageExists(pkg) {
			pack, err := upstream.fetchPackage(pkg)
			if err != nil {
				log.Error(err)
				continue
			}
			log.Infof("Downloading package %s for target %s", pack.Name, pack.Target)
			file := upstream.downloadPackage(pack)
			files = append(files, file)
		} else {
			log.Info(fmt.Sprintf("Dependancy %s exists in target, skipping download", pkgName))
		}
	}

	pack, err := upstream.fetchPackage(p)
	if err != nil {
		log.Error(err)
		return
	}
	pkgName := fmt.Sprintf("%s/%s/%s/%s", p.Origin, p.Name, p.Version, p.Release)
	log.Info(fmt.Sprintf("Downloading package %s for target %s", pack.Name, pack.Target))
	file := upstream.downloadPackage(pack)
	log.Infof("Uploading package %s to channel %s", pkgName, channel)
	packageUpload(target, file, channel)

	// This is a safe guard, sometimes bad things happen on upload where we cannot sync the package to
	// a channel. This will ensure the promotion is atleast attempted.
	packagePromote(target, pkgName, channel, pack.Target)
	files = append(files, file)

	log.Info("Cleaning up downloaded files")
	for _, file := range files {
		log.Debug("Removing file ", file)
		os.Remove(file)
	}
}
