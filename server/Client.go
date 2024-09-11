package server

import (
	"bufio"
)

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
			client = server.Broadcast(client, client.Pseudo, "leave")
			// fmt.Printf(client.Pseudo, " disconnected.\n")
			break
		}
		//Envoie du message à tout les utilisateurs
		client = server.Broadcast(client, message, "message")
	}
}
