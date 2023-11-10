package main

import (
	"fmt"
	"time"
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseFiles("web/index.html"))
var hidden_word []string // Variable sous forme de liste qui contient le mot caché
var used_letters []string // Variable sous forme de liste qui contient les lettres utilisées

type Page struct {
	Valeur string
}

func main() {

	http.HandleFunc("/home", homeHandler)
	http.ListenAndServe(":8080", nil) // Accet pour le port 8080
	http.Handle("/", http.FileServer(http.Dir("static/"))) // A utiliser si notre site est static

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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Valeur: "hagman web"}
	err := templates.ExecuteTemplate(w, "web/index.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}