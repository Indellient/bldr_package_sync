package main

// A Definition of a habitat origin as interpreted from
// configuration files.
type Origin struct {
	Name     string   `toml:"name"`
	Channels []string `toml:"channels"`
}

// Structure of a configuration file.
type Config struct {
	Upstream BldrApi
	Target   BldrApi
	Env      []string
	Interval int
	Origins  []Origin `toml:"origin"`
}
