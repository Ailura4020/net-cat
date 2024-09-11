package server

import (
	"bufio"
	"net"
)

// ?Fonction qui check si le nom entré est déjà pris ou non
func (server *Server) User(conn net.Conn) (bool, string) {
	buf := bufio.NewReader(conn)
	name, _ := buf.ReadString('\n')
	//On boucle sur les pseudos déjà rentrés
	for _, pseudo := range server.clients {
		//On check si le pseudo est vide ou déjà pris
		if string(pseudo.Pseudo) == name[:len(name)-1] || len(name) == 1 {
			conn.Write([]byte("Name already taken, enter a new name: "))
			//Si le pseudo est déjà pris on return false
			return false, name
		}
	}
	//Sinon on return true
	return true, name
}

// ?Fonction qui gère le duplicata de rename
func (server *Server) RenameDeplicates(client Client, newname string) bool {
	//On range sur les pseudos déjà existants
	for _, name := range server.clients {
		//Si le pseudo est déjà pris on return false
		if newname == name.Pseudo {
			return false
		}
	}
	//Sinon on return true
	return true
}
