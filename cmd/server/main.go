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

const (
	IP   = "localhost"
	PORT = "8080"
)

func read(conn net.Conn) {
	message, err := bufio.NewReader(conn).ReadString('\n')
	gestionErreur(err)

	fmt.Print("Client:", string(message))

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
	// var mystring [10]string

	for {
		conn, err := ln.Accept()
		// fmt.Println("LOCAL: ", conn.LocalAddr())
		// fmt.Println("REMOTE: ", conn.RemoteAddr())
		if err == nil {
			clients = append(clients, conn) //quand un client se connecte on le rajoute à notre tableau
		}
		gestionErreur(err)
		fmt.Println("Un client est connecté depuis", conn.RemoteAddr())
		fmt.Println("CLIENTS: ", clients)
		go func() { // création de notre goroutine quand un client est connecté
			// fmt.Print("Enter: ")
			// fmt.Scan(&mystring[0])
			// fmt.Println("Hello", mystring[0])
			buf := bufio.NewReader(conn)

			for {
				name, err := buf.ReadString('\n')

				if err != nil {
					fmt.Printf("Client disconnected.\n")
					break
				}
				for _, c := range clients {
					fmt.Println("C LOCAL: ", c.RemoteAddr())
					// fmt.Println("C LOCAL: ", c.Read())
					// fmt.Scan(&mystring[i])
					// fmt.Printf("Welcome, %v [%v]", i, mystring[i])
					// fmt.Print("BYTE: name: %v msg: %v", name, []byte(name))
					fmt.Println("NAME: ", name)

					c.Write([]byte(name)) // on envoie un message à chaque client
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
