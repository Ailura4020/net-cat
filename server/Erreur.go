package server

import "fmt"

// ?Fonction qui gère les erreurs
func GestionErreur(err error) {
	if err != nil {
		fmt.Println("Error encountered: ", err)
		return
	}
}
