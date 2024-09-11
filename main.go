package main

import (
	"netcat/server"
	"os"
	"fmt"
)

func main() {
	//On defini le port du server
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments, usage: go run . [Port number]")
		return
	}
	port := os.Args[1]

	//On défini le port et l'ip du server
	server := server.Server{
		IP: "localhost",
		PORT: port,
	}
	
	//On lance le server
	server.Run()
}
