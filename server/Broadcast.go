package server

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