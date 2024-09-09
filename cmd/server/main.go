package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	conn   net.Conn
	Pseudo string
	Messages []string
}

type Identification struct {
	Name []string
	Id   []interface{}			
}

const (
	IP   = "localhost"
	PORT = "8081"
)

func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
}

func read(conn net.Conn) {
	message, err := bufio.NewReader(conn).ReadString('\n')
	gestionErreur(err)
	fmt.Print("Client: ", string(message))
}

func main() {
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
	
	//Déclaration des structures
	// var clients []net.Conn // tableau de clients
	// var tab Identification
	
	for {
		//Autorisation d'une nouvelle connection
		conn, err := ln.Accept()
		//Gestion d'erreur
		if err != nil {
			fmt.Println(err)
			continue
		}
		client := Client{
			conn: conn,
		}
		// clients = append(clients, conn) //quand un client se connecte on le rajoute à notre tableau
		//Ajout de l'adresse IP dans notre structure
		// tab.Id = append(tab.Id, conn.RemoteAddr())
		fmt.Println("Remote Addr: ", conn.RemoteAddr())
		
		//Demande du nom
		client.conn.Write([]byte("Welcome\n"))
		client.conn.Write([]byte("Enter your name: "))

		//Vérification si le nom choisi est déjà pris
		name := client.User(conn)
		
		//Ajout du nom au tableau de noms
		messages := []string{}
		client = Client{
			conn: conn,
			Pseudo: name[:len(name)-1],
			Messages: messages,
		}
		// clients.pseudo = append(clients.Pseudo, name[:len(name)-1])

		// création de notre goroutine quand un client est connecté
		go client.HandleConnection(client)
		fmt.Println("test3")
	}
}

func (clients *Client) HandleConnection(client Client) {
	// Close the connection when we're done
	// defer client.conn.Close()
	fmt.Println("test1")
	buf := bufio.NewReader(client.conn)
	for {
		message, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("test2")
			fmt.Printf("Client disconnected.\n")
			break
		}
		// client.Messages = append(client.Messages, message)
		fmt.Println("test5")
		client.conn.Write([]byte("[" + string(client.Pseudo) + "]: "))
		fmt.Println("test6")
		client.conn.Write([]byte(message)) // on envoie un message à chaque client
		fmt.Println("test4")
	}
	// conn.Write([]byte(buf))
}

func (clients *Client) User(conn net.Conn) string {
	buf := bufio.NewReader(clients.conn)
	name, _ := buf.ReadString('\n')
	for _, pseudo := range clients.Pseudo {
		if string(pseudo) == name[:len(name)-1] || len(name) == 1 {
			conn.Write([]byte("Enter a new name: "))
			clients.User(conn)
		}
	}
	return name
}
