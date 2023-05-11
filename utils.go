package main

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var config Configuration

func init() {
	config.Address = "0.0.0.0:8080"
	config.ReadTimeout = 10
	config.WriteTimeout = 600
	config.Static = "public"
}

