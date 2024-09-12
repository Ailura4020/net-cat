package server

import (
	"fmt"
	"net"
)

// ?Fonction qui va lancer le server et attribuer des goroutines aux utilisateurs
func (server *Server) Run() {
	//Message au lancement
	fmt.Println("Lancement du serveur...")

	//Création d'une connection au port et à l'Ip donnée
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.IP, server.PORT))
	GestionErreur(err)

	client := Client{}
	go server.AdminConnection(client)

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

			connected = ""
			for _, connect := range server.clients {
				connected += connect.Pseudo + ", "
			}

			//Message de bienvenue
			client.conn.Write([]byte("Welcome\n"))

			//Affichage des clients déjà connectés
			if len(connected) == 0 {
				//Si aucun client n'est connecté
				client.conn.Write([]byte("Server empty\n"))
			} else {
				//Sinon on affiche tout les clients
				client.conn.Write([]byte("Clients connected: " + "\033[34m" + connected[:len(connected)-2] + "\033[0m" + "\n"))
			}

			//Demande du nom
			client.conn.Write([]byte("Enter your name: "))

			//Vérification si le nom choisi est déjà pris
			duplicate, name := server.User(conn)
			//Tant que la variable booleene est false (duplicate), on redemande en boucle un pseudo
			for !duplicate {
				duplicate, name = server.User(conn)
			}

			//Affichage de l'arrivé d'un client aux autres utilisateurs
			client = server.Broadcast(client, name[:len(name)-1], "join")

			//Ajout du nom au tableau de noms
			client = Client{
				conn:   conn,
				Pseudo: name[:len(name)-1],
			}

			//Notifie le server quand un client se connecte
			// fmt.Printf(client.Pseudo, " connected.\n")
			fmt.Println("Number of clients connected: ", len(server.clients))

			//Ajout de la structure client à la structure server
			server.mutex.Lock()
			server.clients = append(server.clients, client)
			server.mutex.Unlock()

			//Création de notre goroutine quand un client est connecté
			go server.HandleConnection(client)
			//(La goroutine va executer une fonction en parrallèle de notre server)
		}
	}
}
