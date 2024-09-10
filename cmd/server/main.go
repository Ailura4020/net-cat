package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Server struct {
	clients []Client
}

type Client struct {
	conn   net.Conn
	Pseudo string
	// Messages []string
}

type Identification struct {
	Name []string
	Id   []interface{}			
}

const (
	IP   = "localhost"
	PORT = "8081"
)

func main() {
	server := Server{
	}
	server.Run()
}

//Fonction qui gère les erreurs
func gestionErreur(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

//Fonction qui va lancer le server et attribuer des goroutines aux utilisateurs
func (server *Server) Run() {
	//Message au lancement
	fmt.Println("Lancement du serveur...")

	//On defini le port du server
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments, usage: go run . [port number]")
		return
	}
	port := os.Args[1]

	//Création d'une connection au port et à l'Ip donnée
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, port))
	gestionErreur(err)
	// messages := []string{}
	
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
			// ascii := AsciiArt()
			//Demande du nom
			// client.conn.Write([]byte(ascii))
			client.conn.Write([]byte("Welcome\n"))
			client.conn.Write([]byte("Enter your name: "))
			
			//Vérification si le nom choisi est déjà pris
			name := client.User(conn)
					
			//Affichage de l'arrivé d'un client aux autres utilisateurs
			server.Broadcast(client, name[:len(name)-1], 0)
					
			//Ajout du nom au tableau de noms
			client = Client{
				conn: conn,
				Pseudo: name[:len(name)-1],
				// Messages: messages,
			}
			fmt.Println(len(server.clients))
			// fmt.Println(client.Messages)
			
			//Ajout de la structure client à la structure server
			server.clients = append(server.clients, client)
					
			// création de notre goroutine quand un client est connecté
			go server.HandleConnection(client)
		}
	}
}

//Fonction qui gère l'envoie des messages des utilisateurs
func (server *Server) HandleConnection(client Client) {
	// Close the connection when we're done
	// defer client.conn.Close()

	// for _, name := range server.clients {
	// 	for _, historic := range client.Messages {
	// 		name.conn.Write([]byte(historic))
	// 	}
	// }
	buf := bufio.NewReader(client.conn)
	for {
		message, err := buf.ReadString('\n')
		if err != nil {
			server.Broadcast(client, client.Pseudo, 1)
			fmt.Printf("Client disconnected.\n")
			break
		}
		//Envoie du message à tout les utilisateurs
		server.Broadcast(client, message, 2)
	}
	// conn.Write([]byte(buf))
}

//Fonction qui check si le nom entré est déjà pris ou non
func (clients *Client) User(conn net.Conn) string {
	buf := bufio.NewReader(clients.conn)
	name, _ := buf.ReadString('\n')
	//On boucle sur les pseudos déjà rentrés
	for _, pseudo := range clients.Pseudo {
		//On check si le pseudo est vide ou déjà pris
		if string(pseudo) == name[:len(name)-1] || len(name) == 1 {
			//Si c'est le cas on utilise la récursivité pour redemander le pseudo
			conn.Write([]byte("Name already taken, enter a new name: "))
			clients.User(conn)
		}
	}
	return name
}

//Fonction qui envoie le message à tout les utilisateurs
func (server *Server) Broadcast(client Client, message string, messagetype int) {
	// fmt.Println("Pseudo: ", client.Pseudo)
	if messagetype == 0 {
		for _, name := range server.clients {
			name.conn.Write([]byte(message + " has joined the chat.\n"))
		}
	} else if messagetype == 1 {
			for i, name := range server.clients {
				name.conn.Write([]byte(message + " has left the chat.\n"))
				if name == client {
					server.clients = append(server.clients[:i], server.clients[i+1:]...)
					fmt.Println(server.clients)
				}
			}
		} else if messagetype == 2 {
			// client.Messages = append(client.Messages, message)
			for _, name := range server.clients {
			name.conn.Write([]byte("[" + string(client.Pseudo) + "]: "))
			name.conn.Write([]byte(message))
			// fmt.Println(name)
		}
	}
}

func AsciiArt() string {
	return `
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    ".       | "' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     "-'       '--'
	` + "\n"
}
