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
	// // Listen for incoming connections on port 8080
	// ln, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// // Accept incoming connections and handle them
	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}

	// 	// Handle the connection in a new goroutine
	// 	go handleConnection(conn)
	// }
	fmt.Println("Lancement du serveur ...")

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	gestionErreur(err)

	var clients []net.Conn // tableau de clients
	var tab Identification

	for {
		conn, err := ln.Accept()
		if err == nil {
			clients = append(clients, conn) //quand un client se connecte on le rajoute à notre tableau
			tab.Id = append(tab.Id, conn.RemoteAddr())
			// id.Id = tab
		}
		buf := bufio.NewReader(conn)
		conn.Write([]byte("Enter your name: "))
		name, err := buf.ReadString('\n')
		if tab.User(name) {
			fmt.Println(tab.User(name))
			tab.Name = append(tab.Name, name[:len(name)-1])
		} else {
			buf := bufio.NewReader(conn)
			conn.Write([]byte("Enter your name: "))
			name, _ := buf.ReadString('\n')
			fmt.Println(name)
		}
		fmt.Println(tab)
		// fmt.Println(tab.User(name))
		// conn.Write([]byte(tab.Name))
		// conn.Write([]byte(" has joined the chat."))

		gestionErreur(err)
		// fmt.Println(tab)
		// fmt.Println("Un client est connecté depuis", conn.RemoteAddr())
		// fmt.Println("CLIENTS: ", clients)

		go func() { // création de notre goroutine quand un client est connecté
			// fmt.Println("Hello", mystring[0])
			buf := bufio.NewReader(conn)
			for {
				for i, c := range clients {
					message, err := buf.ReadString('\n')
					if err != nil {
						fmt.Printf("Client disconnected.\n")
						break
					}
					c.Write([]byte("[" + tab.Name[i] + "]: "))
					c.Write([]byte(message)) // on envoie un message à chaque client
				}
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

func (tab *Identification) User() {
	buf := bufio.NewReader(conn)
	conn.Write([]byte("Enter your name: "))
	name, _ := buf.ReadString('\n')
	for _, pseudo := range tab.Name {
		if pseudo == name {
			tab.User()
		}
	}
	// conn.Write([]byte("Enter your name: "))
	// name, err := buf.ReadString('\n')
	tab.Name = append(tab.Name, name[:len(name)-1])
}
