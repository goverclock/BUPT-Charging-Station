package main

import (
	"buptcs/data"
	"strconv"
)

const (
	CodeSucceed   int = 200
	CodeKeyError  int = 400
	CodeForbidden int = 403
)

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var config Configuration

func init() {
	addr := "0.0.0.0"
	config.Address = addr + ":" + strconv.Itoa(data.Port)
	config.ReadTimeout = 10
	config.WriteTimeout = 600
	config.Static = "public"
}
