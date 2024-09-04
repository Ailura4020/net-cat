package main

import (
	"bufio"
	"fmt"
	"net"
)

func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
}

type Client struct {
	conn   net.Conn
	pseudo string
}

type Identification struct {
	Name []string
	Id   []interface{}
}

const (
	IP   = "localhost"
	PORT = "8081"
)

func read(conn net.Conn) {
	message, err := bufio.NewReader(conn).ReadString('\n')
	gestionErreur(err)

	fmt.Print("Client: ", string(message))

}

func main() {
	fmt.Println("Lancement du serveur ...")

	//Création d'une connection au port et à l'Ip donnée
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	gestionErreur(err)
	
	//Déclaration des structures
	var clients []net.Conn // tableau de clients
	var tab Identification
	
	for {
		//Autorisation d'une nouvelle connection
		conn, err := ln.Accept()
		//Gestion d'erreur
		if err != nil {
			fmt.Println(err)
			continue
		}
		
		// clients = append(clients, conn) //quand un client se connecte on le rajoute à notre tableau
		//Ajout de l'adresse IP dans notre structure
		tab.Id = append(tab.Id, conn.RemoteAddr())
		
		//Demande du nom
		conn.Write([]byte("Welcome\n"))
		conn.Write([]byte("Enter your name: "))

		//Vérification si le nom choisi est déjà pris
		name := tab.User(conn)
		//Ajout du nom au tableau de noms
		tab.Name = append(tab.Name, name[:len(name)-1])

		// gestionErreur(err)
		// fmt.Println("Un client est connecté depuis", conn.RemoteAddr())
		// fmt.Println("CLIENTS: ", clients)

		// création de notre goroutine quand un client est connecté
		fmt.Println("Netconn: ", clients)
		go func() { 
			// buf := bufio.NewReader(conn)
			for {
				// for _, c := range clients {
				// 	message, err := buf.ReadString('\n')
				// 	if err != nil {
				// 		fmt.Printf("Client disconnected.\n")
				// 		break
				// 	}
				// 	fmt.Println(c)
				// 	c.Write([]byte("[" + c + "]: "))
				// 	c.Write([]byte(message)) // on envoie un message à chaque client
				// }
			}
		}()
	}
}

func handleConnection(conn net.Conn) {
	// Close the connection when we're done
	defer conn.Close()

	// Read incoming data
	buf := make([]byte, 1024)
	// buf := make([]byte, 4096)
	// bufw := make([]byte, 1024)
	_, err := conn.Read(buf)
	// _, err2 := conn.Write(buf)
	// fmt.Println(test)

	fmt.Println("LOCAL: ", conn.LocalAddr())
	fmt.Println("REMOTE: ", conn.RemoteAddr())
	fmt.Println(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the incoming data
	// fmt.Printf("Received: %s", buf)
	conn.Write([]byte(buf))
}

func (tab *Identification) User(conn net.Conn) string {
	buf := bufio.NewReader(conn)
	name, _ := buf.ReadString('\n')
	for _, pseudo := range tab.Name {
		if pseudo == name[:len(name)-1] {
			conn.Write([]byte("Enter a new name: "))
			tab.User(conn)
		}
	}
	return name
}
