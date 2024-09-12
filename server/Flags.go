package server

import (
	"strings"
	"time"
)

func (server *Server) FlagRename(client Client, message string) Client {
	rename := false
	index := 0
	//On prend ce qu'il y'a derrière le /rename
	newname := strings.TrimSpace(message[7:])
	//On range sur les clients pour checker les pseudo deja pris un à un
	for i, name := range server.clients {
		if server.RenameDeplicates(client, newname) {
			name.conn.Write([]byte(string(client.Pseudo) + " has changed his/her name for: " + newname + "\n"))
			//On check si le nouveau nom choisi est déjà pris
			if name.conn == client.conn {
				rename = true
				index = i
			}
		}
	}
	Txt = append(Txt, "["+time.Now().Format("2006-01-02 15:04:05")+string(server.clients[index].Pseudo)+"]"+" has changed his/her name for: "+newname+"\n")
	if rename && server.RenameDeplicates(client, newname) {
		//On modifie la structure donc on lock avec le mutex le temps des changements
		server.mutex.Lock()
		//On change la structure client
		client.Pseudo = newname
		//On change aussi la structure du client stockée dans la structure server
		server.clients[index].Pseudo = newname
		server.mutex.Unlock()
	} else if !server.RenameDeplicates(client, newname) {
		//On affiche un message à l'utilisateur qui utilise /rename si le rename n'est pas possible
		client.conn.Write([]byte("Name already taken, choose another one.\n"))
	}
	return client
}

func (server *Server) FlagColor(client Client, index int, message string) Client {
	//On prend ce qu'il y'a derrière le /rename
	newname := strings.ToLower(strings.TrimSpace(message[6:]))

	//Booleen qui va vérifier si la couleur est bonne ou non
	color := true

	//On modifie la structure donc on lock avec le mutex le temps des changements
	server.mutex.Lock()

	//Modifier la couleur du pseudo en fonction de la couleur choisie
	switch newname {
	case "yellow":
		client.Pseudo = "\033[33m" + server.clients[index].Pseudo + "\033[0m"
	case "red":
		client.Pseudo = "\033[31m" + server.clients[index].Pseudo + "\033[0m"
	case "blue":
		client.Pseudo = "\033[34m" + server.clients[index].Pseudo + "\033[0m"
	case "magenta":
		client.Pseudo = "\033[35m" + server.clients[index].Pseudo + "\033[0m"
	case "cyan":
		client.Pseudo = "\033[36m" + server.clients[index].Pseudo + "\033[0m"
	case "green":
		client.Pseudo = "\033[32m" + server.clients[index].Pseudo + "\033[0m"
	case "white":
		client.Pseudo = "\033[97m" + server.clients[index].Pseudo + "\033[0m"
	default:
		//Si le client entre une couleur non valide
		client.conn.Write([]byte("Invalid color, choose another one\n"))
		color = false
	}
	//Unlock du mutex une fois les changements terminés
	server.mutex.Unlock()

	//Affichage aux autres clients si la couleur est bonne
	if color {
		//Affichage de la modification à tout les autres clients
		for _, name := range server.clients {
			name.conn.Write([]byte(string(server.clients[index].Pseudo) + " has changed his/her color for: " + newname + "\n"))
		}
		//Stockage du changement dans le tableau qui va gérer les logs
		Txt = append(Txt, "["+time.Now().Format("2006-01-02 15:04:05")+string(server.clients[index].Pseudo)+" has changed his/her color for: "+newname+"\n")
	}
	return client
}
