package main

import "buptcs/server"

func main() {
	s := server.New()
	s.Run(":8080")
}
