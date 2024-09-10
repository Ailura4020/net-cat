package server

// ?Fonction qui check si le nom entré est déjà pris ou non
func (clients *Client) User(conn net.Conn) string {
	buf := bufio.NewReader(clients.conn)
	name, _ := buf.ReadString('\n')
	//On boucle sur les pseudos déjà rentrés
	for _, pseudo := range clients.Pseudo {
		//On check si le pseudo est vide ou déjà pris
		if string(pseudo) == name[:len(name)-1] || len(name) == 1 {
			//Si c'est le cas on utilise la récursivité pour redemander le pseudo
			conn.Write([]byte("Name already taken, enter a new name: "))
			clients.User(conn)
		}
	}
	return name
}