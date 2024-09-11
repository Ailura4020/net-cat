package main

import (
	"fmt"
	"log"
	"netcat/server"
	"os"
)

func main() {
	//On defini le port du server
	if len(os.Args) < 2 {
		fmt.Println("Wrong number of arguments, usage: go run . [Port number] [IP adress]")
		return
	}
	port := os.Args[1]
	ip := ""
	if len(os.Args) == 2 {
		ip = "localhost"
	} else {
		ip = os.Args[2]
	}

	//On défini le port et l'ip du server
	server := server.Server{
		IP:   ip,
		PORT: port,
	}
	//Permet d'effacer l'historique dans log.txt
	file, err := os.OpenFile("../net-cat/log.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	//On lance le server
	server.Run()
}
