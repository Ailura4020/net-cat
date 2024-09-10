package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Server struct {
	clients []Client
	mutex   sync.Mutex
}

type Client struct {
	conn   net.Conn
	Pseudo string
	// Messages []string
}

var Log []Historic

type Historic struct {
	Time    string
	Pseudo  string
	Message string
}

const (
	IP   = "localhost"
	PORT = "8081"
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

// ?Fonction qui gère l'envoie des messages des utilisateurs
func (server *Server) HandleConnection(client Client) {
	// Close the connection when we're done
	// defer client.conn.Close()

	//Affichage de l'historique des messages
	for _, historic := range Log {
		client.conn.Write([]byte("\033[33m" + "[" + historic.Time + "]" + "[" + historic.Pseudo + "]: " + historic.Message + "\033[0m"))
	}

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
}

// ?Fonction qui check si le nom entré est déjà pris ou non
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

// fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), color.Magenta.Sprintf(sender.pseudo), color.Red.Sprintf(message)

// ?Fonction qui envoie le message à tout les utilisateurs
func (server *Server) Broadcast(client Client, message string, messagetype int) {
	if messagetype == 0 {
		for _, name := range server.clients {
			name.conn.Write([]byte("\033[32m" + time.Now().Format("2006-01-02 15:04:05") + "] " + message + " has joined the chat.\n" + "\033[0m"))
		}
	} else if messagetype == 1 {
		for i, name := range server.clients {
			name.conn.Write([]byte("\033[31m" + "[" + time.Now().Format("2006-01-02 15:04:05") + "] " + message + " has left the chat.\n" + "\033[0m"))
			if name == client {
				server.clients = append(server.clients[:i], server.clients[i+1:]...)
				fmt.Println(server.clients)
			}
		}
	} else if messagetype == 2 {
		//Enregistrement des informations du message dans le tableau de logs
		historic := Historic{
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Pseudo:  client.Pseudo,
			Message: message,
		}

		server.mutex.Lock()
		Log = append(Log, historic)
		server.mutex.Unlock()

		if strings.HasPrefix(message, "/rename") {
			newname := strings.Split(message, " ")
			for i, name := range server.clients {
				name.conn.Write([]byte(string(client.Pseudo) + " has changed his name for: " + newname[1]))
				if name == client {
					// fmt.Println("Before changes: ", server.clients[i].Pseudo)
					// server.clients[i].Pseudo = newname[1][:len(newname[1])-1]
					// fmt.Println("After changes: ", server.clients[i].Pseudo)
					client := Client{
						conn:   client.conn,
						Pseudo: newname[1][:len(newname[1])-1],
					}
					server.clients[i] = client
				}
			}
		} else {
			for _, name := range server.clients {
				name.conn.Write([]byte("\033[37m" + "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "\033[36m" + "[" + string(client.Pseudo) + "]: " + "\033[0m"))
				name.conn.Write([]byte(message))
			}
		}
	}
}

// ?Fonction qui affiche un pingoin
func AsciiArt() string {
	return `
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
   FqM            MMMM
 __| ".        |\dS"qML
 |    ".       | "' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     "-'       '--'
	` + "\n"
}
