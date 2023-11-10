package main

import (
	"fmt"
	"time"
	"net/http"
	"text/template"
)

const port = ":8080"
var hidden_word []string // Variable sous forme de liste qui contient le mot caché
var used_letters []string // Variable sous forme de liste qui contient les lettres utilisées

func home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home")
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("./templates/" + tmpl + ".page.tmpl")
	if err != nil {
		fmt.Println("error")
	}
	t.Execute(w, nil)
}

func main() {

	http.HandleFunc("/", home)

	PrintHangmanAscii() // Appel de la fonction pour afficher "Hangman" en ascii
	fmt.Println("Bienvenue sur Hangman !")
	fmt.Println("La partie va bientôt commencer..")
	time.Sleep(4 * time.Second)
	Clear()
	tries := 0 // Initialisation de la variable des essais
	originalWord := WordPicker(RandomNumber()) // Initalisation du mot aléatoire a faire deviner
	Hidden(originalWord) // Modificiton du mot généré en underscore
	RunHangman(originalWord, tries) // Lancement du jeu

	fmt.Println("(https://localhost:8080) - Serveur started on port", port)
	http.ListenAndServe(port, nil)
}