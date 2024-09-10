package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"netcat/server"
	"sync"
)

// ?Fonction qui va lancer le server et attribuer des goroutines aux utilisateurs
func (server *Server) Run() {
	//Message au lancement
	fmt.Println("Lancement du serveur...")

	//On defini le port du server
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments, usage: go run . [Port number]")
		return
	}
	port := os.Args[1]

	//Création d'une connection au port et à l'Ip donnée
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, port))
	gestionErreur(err)

	for {
		//Autorisation d'une nouvelle connection
		conn, err := ln.Accept()
		//Gestion d'erreur
		if err != nil {
			fmt.Println(err)
			continue
		}

		//Ajout d'une variable qui va stocker les données de la nouvelle connection
		client := Client{
			conn: conn,
		}

		//Verification du nombre de clients déjà connectés
		if len(server.clients) == 10 {
			client.conn.Write([]byte("Server is full, 10 Users already connected.\n"))
			client = Client{}
		} else {
			//Affichage du logo linux
			ascii := AsciiArt()
			client.conn.Write([]byte(ascii))

			//Demande du nom
			client.conn.Write([]byte("Welcome\n"))
			client.conn.Write([]byte("Enter your name: "))

			//Vérification si le nom choisi est déjà pris
			name := client.User(conn)

			//Affichage de l'arrivé d'un client aux autres utilisateurs
			server.Broadcast(client, name[:len(name)-1], 0)

			//Ajout du nom au tableau de noms
			client = Client{
				conn:   conn,
				Pseudo: name[:len(name)-1],
			}
			fmt.Println(len(server.clients))
			// fmt.Println(client.Messages)

			//Ajout de la structure client à la structure server
			server.mutex.Lock()
			server.clients = append(server.clients, client)
			server.mutex.Unlock()

			// création de notre goroutine quand un client est connecté
			go server.HandleConnection(client)
		}
	}
}