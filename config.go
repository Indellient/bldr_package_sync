package main

type Origin struct {
	Name     string   `toml:"name"`
	Channels []string `toml:"channels"`
}

type Config struct {
	Upstream BldrApi
	Target   BldrApi
	Interval int
	Origins  []Origin `toml:"origin"`
}
