package main

import (
	// "github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type Syncer struct {
	config Config
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
	return nil
}
