package server

import (
	"fmt"
	"strings"
	"time"
)

// ?Fonction qui envoie le message à tout les utilisateurs
func (server *Server) Broadcast(client Client, message string, messagetype string) Client {
	if messagetype == "join" {
		//Enregistrement des informations du message dans le tableau de logs
		historic := Historic{
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Pseudo:  message,
			Message: "has joined the chat.\n",
		}
				
		//Lock des autres clients le temps de changer la base de donnée
		server.mutex.Lock()
		Log = append(Log, historic)
		server.mutex.Unlock()
			
		for _, name := range server.clients {
			name.conn.Write([]byte("\033[32m" + time.Now().Format("2006-01-02 15:04:05") + "] " + message + " has joined the chat.\n" + "\033[0m"))
		}
	} else if messagetype == "leave" {
		//Enregistrement des informations du message dans le tableau de logs
		historic := Historic{
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Pseudo:  message,
			Message: "has left the chat.\n",
		}

		//Lock des autres clients le temps de changer la base de donnée
		server.mutex.Lock()
		Log = append(Log, historic)
		server.mutex.Unlock()

		for i, name := range server.clients {
			name.conn.Write([]byte("\033[31m" + "[" + time.Now().Format("2006-01-02 15:04:05") + "] " + message + " has left the chat.\n" + "\033[0m"))
			if name == client {
				server.clients = append(server.clients[:i], server.clients[i+1:]...)
				fmt.Println(server.clients)
			}
		}
	} else if messagetype == "message" {
		//Enregistrement des informations du message dans le tableau de logs
		historic := Historic{
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Pseudo:  client.Pseudo,
			Message: message,
		}

		//Lock des autres clients le temps de changer la base de donnée
		server.mutex.Lock()
		Log = append(Log, historic)
		server.mutex.Unlock()

		//Filtrer si le message est un rename ou pas
		if strings.HasPrefix(message, "/rename") && strings.TrimSpace(message) != "/rename" {
			rename := false
			index := 0
			//On prend ce qu'il y'a derrière le /rename
			newname := strings.TrimSpace(message[7:])
			//On range sur les clients pour checker les pseudo deja pris un à un
			for i, name := range server.clients {
				if server.RenameDeplicates(client ,newname) {
					name.conn.Write([]byte(string(client.Pseudo) + " has changed his/her name for: " + newname + "\n"))
					//On check si le nouveau nom choisi est déjà pris
					if name == client {
						rename = true
						index = i
					}
				}
			}
			if rename && server.RenameDeplicates(client, newname) {
				//On modifie la structure donc on lock avec le mutex le temps des changements
				server.mutex.Lock()
				//On change la structure client
				client.Pseudo = newname
				//On change aussi la structure du client stockée dans la structure server
				server.clients[index].Pseudo = newname
				server.mutex.Unlock()
			} else if !server.RenameDeplicates(client ,newname) {
				//On affiche un message à l'utilisateur qui utilise /rename si le rename n'est pas possible
				client.conn.Write([]byte("Name already taken, choose another one.\n"))
			}
			rename = false
		} else {
			//Si le message n'est pas rename, on affiche juste le message à tout les clients
			for _, name := range server.clients {
				name.conn.Write([]byte("\033[37m" + "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "\033[36m" + "[" + string(client.Pseudo) + "]: " + "\033[0m"))
				name.conn.Write([]byte(message))
			}
		}
	}
	//On return la structure client, modifiée ou non
	return client
}
