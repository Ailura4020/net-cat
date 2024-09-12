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
		Txt = append(Txt, ("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + message + " has joined the chat.\n"))
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
		Txt = append(Txt, (time.Now().Format("2006-01-02 15:04:05") + "] " + message + " has left the chat.\n"))
		server.mutex.Unlock()

		for i, name := range server.clients {
			name.conn.Write([]byte("\033[31m" + "[" + time.Now().Format("2006-01-02 15:04:05") + "] " + message + " has left the chat.\n" + "\033[0m"))
			if name == client {
				server.clients = append(server.clients[:i], server.clients[i+1:]...)
				fmt.Println(server.clients)
			}
		}

		//Check si le message est envoyé par l'admin
	} else if messagetype == "admin" {
		//Enregistrement des informations du message dans le tableau de logs
		historic := Historic{
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Pseudo:  "Admin",
			Message: message,
		}

		//Lock des autres clients le temps de changer la base de donnée
		server.mutex.Lock()
		Log = append(Log, historic)
		Txt = append(Txt, ("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + "[" + "Admin" + "]" + message))
		server.mutex.Unlock()

		for _, name := range server.clients {
			name.conn.Write([]byte("\033[31m" + "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + "Admin" + "]: " + message + "\033[0m" + "\n"))
		}
	} else if messagetype == "message" {
		index := 0
		//On recupère l'index
		for i, name := range server.clients {
			if name.conn == client.conn {
				index = i
			}
		}

		//Enregistrement des informations du message dans le tableau de logs
		historic := Historic{
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Pseudo:  server.clients[index].Pseudo,
			Message: message,
		}

		//Lock des autres clients le temps de changer la base de donnée
		server.mutex.Lock()
		Log = append(Log, historic)
		Txt = append(Txt, "["+time.Now().Format("2006-01-02 15:04:05")+"] "+"["+client.Pseudo+"]"+message)
		server.mutex.Unlock()

		//TODO Filtrer si le message est un rename ou pas
		if strings.HasPrefix(message, "/rename") && strings.TrimSpace(message) != "/rename" {
			client = server.FlagRename(client, message)

			//TODO Gestion si l'utilisateur cherche à changer la couleur de son pseudo
		} else if strings.HasPrefix(message, "/color") && strings.TrimSpace(message) != "/color" {
			client = server.FlagColor(client, index, message)
		} else {
			//Si le message n'est pas rename, on affiche juste le message à tout les clients
			for _, name := range server.clients {
				name.conn.Write([]byte("\033[37m" + "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + "\033[36m" + string(client.Pseudo) + "\033[0m" + "]:"))
				name.conn.Write([]byte(message))
			}
		}
	}
	//Appelle la fonction pour imprimer l'historique en fichier txt
	LogHistory()
	//On return la structure client, modifiée ou non
	return client
}
