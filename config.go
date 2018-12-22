package main

type Config struct {
	Upstream BldrApi
	Target   BldrApi
	Interval int
	Origins  []string
}
