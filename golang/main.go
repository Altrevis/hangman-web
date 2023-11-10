package main

import (
	"fmt"
	"time"
	"net/http"
)


var hidden_word []string // Variable sous forme de liste qui contient le mot caché
var used_letters []string // Variable sous forme de liste qui contient les lettres utilisées

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Bienvenue sur la page d'accueil!</h1>")
}

func main() {

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

	PrintHangmanAscii() // Appel de la fonction pour afficher "Hangman" en ascii
	fmt.Println("Bienvenue sur Hangman !")
	fmt.Println("La partie va bientôt commencer..")
	time.Sleep(4 * time.Second)
	Clear()
	tries := 0 // Initialisation de la variable des essais
	originalWord := WordPicker(RandomNumber()) // Initalisation du mot aléatoire a faire deviner
	Hidden(originalWord) // Modificiton du mot généré en underscore
	RunHangman(originalWord, tries) // Lancement du jeu
}