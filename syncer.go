package main

import (
	// "github.com/BurntSushi/toml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	// "sync"
)

type Syncer struct {
	config Config
}

func (syncer Syncer) syncPackages(origin string, channel string, upstream BldrApi, target BldrApi) bool {
	log.Debug("Beginning the package sync process")

	upstreamPkgs := upstream.listAllPackages(origin, channel)
	log.Debug(fmt.Sprintf("Found %d packages on %s", len(upstreamPkgs.Data), upstream.Url))

	targetPkgs := target.listAllPackages(origin, channel)
	log.Debug(fmt.Sprintf("Found %d packages on %s", len(targetPkgs.Data), target.Url))

	// Good enough to figure out the difference before calculating deps
	pkgDatas := packageDifference(upstreamPkgs.Data, targetPkgs.Data)

	log.Debug(fmt.Sprintf("Determining TDEPS for %d packages", len(pkgDatas)))

	// var wg sync.WaitGroup
	for _, p := range pkgDatas {
		// wg.Add(1)
		// go func(wg *sync.WaitGroup) {
		deps := upstream.fetchPackageDeps(p)
		for _, pkg := range deps {
			pkgName := fmt.Sprintf("%s/%s/%s/%s", pkg.Origin, pkg.Name, pkg.Version, pkg.Release)
			pack := upstream.fetchPackage(pkg)
			// if targPack.Name != "" {
			if !target.packageExists(pack.Ident) {
				log.Debug(fmt.Sprintf("Downloading package %s for target %s", pack.Name, pack.Target))
				file := upstream.downloadPackage(pack)
				log.Debug("Uploading package " + pkgName)
				packageUpload(target, file, "stable")
				os.Remove(file)
			} else {
				log.Debug("Package exists in target " + pkgName)
			}
		}

		pack := upstream.fetchPackage(p)
		pkgName := fmt.Sprintf("%s/%s/%s/%s", p.Origin, p.Name, p.Version, p.Release)
		if !target.packageExists(pack.Ident) {
			log.Debug(fmt.Sprintf("Downloading package %s for target %s", pack.Name, pack.Target))
			file := upstream.downloadPackage(pack)
			log.Debug("Uploading package " + pkgName)
			packageUpload(target, file, "stable")
			os.Remove(file)
		} else {
			log.Debug("Package exists in target " + pkgName)
		}

		// wg.Done()
		// }(&wg)
	}

	// wg.Wait()
	return true
}

func (syncer Syncer) syncKeys(origin string, upstream BldrApi, target BldrApi) bool {
	log.Debug("Beginning the key sync process")
	upstreamKeys := upstream.fetchKeyPaths(origin)
	targetKeys := target.fetchKeyPaths(origin)
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
	for _, origin := range syncer.config.Origins {
		syncer.syncKeys(origin, syncer.config.Upstream, syncer.config.Target)
	}

	for _, origin := range syncer.config.Origins {
		syncer.syncPackages(origin, "stable", syncer.config.Upstream, syncer.config.Target)
	}

	return nil
}
