package server

import (
	"log"
	"os"
)

// ?Fonction qui permet d'obtenir un historique en fichier txt
func LogHistory() {
	//boucle qui imprime le resultat en string
	result := ""
	for _, print := range Txt {
		result += string(print + "\n")
	}

	//permet d'ouvrir le fichier : Create = créer fichier non existant : Wronly = écrire dans le fichier
	file, err := os.OpenFile("../net-cat/log.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//écrire dans le fichier
	_, err = file.WriteString(result)
	if err != nil {
		panic(err)
	}
}
