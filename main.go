package main

import (
	"netcat/server"
)

func main() {
	server := server.Server{}
	server.Run()
}
