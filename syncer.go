package main

import (
	// "github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type Syncer struct {
	config Config
}

func (syncer Syncer) syncKeys() bool {
	log.Debug("Beginning the key sync process")
	return true
}

func (syncer Syncer) run() error {
	syncer.syncKeys()
	return nil
}
