package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"netcat/server"
	"strings"
	"sync"
	"time"
)

func main() {
	server := Server{}
	server.Run()
}

// ?Fonction qui gère les erreurs
func gestionErreur(err error) {
	if err != nil {
		fmt.Println("Error encountered: ", err)
		return
	}
}
