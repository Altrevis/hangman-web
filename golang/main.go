package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"bufio"
)

// Variables globales pour suivre l'état du jeu
var (
	guesses      = make(map[string]bool) // Stocke les lettres devinées
	wordToGuess  string                   // Mot à deviner
	errorCount   int                      // Nombre d'erreurs
	hangmanState int                      // État actuel du pendu (0 à 6)
)

// Images du pendu pour chaque état
var hangmanStates = []string{
	`
	    +---+
	    |   |
		|
		|
		|
		|
	==========`,
	`
	    +---+
	    |   |
	    O   |
		|
		|
		|
	==========`,
	`
	    +---+
	    |   |
	    O   |
            |   |
		|
		|
	==========`,
	`
	    +---+
	    |   |
	    O   |
   	   /|   |
		|
		|
	==========`,
	`
	    +---+
	    |   |
	    O   |
   	   /|\  |
		|
		|
	==========`,
	`
	    +---+
	    |   |
	    O   |
   	   /|\  |
   	   /    |
		|
	==========`,
	`
	    +---+
	    |   |
	    O   |
           /|\  |
           / \  |
		|
	==========`,
}

// Fonction pour afficher le template HTML
func renderTemplate(w http.ResponseWriter, message string) {
	wordWithState := make([]struct {
		Letter  string
		Guessed bool
	}, len(wordToGuess))

	// Remplir le tableau wordWithState avec les lettres devinées et non devinées
	for i, letter := range wordToGuess {
		wordWithState[i] = struct {
			Letter  string
			Guessed bool
		}{Letter: string(letter), Guessed: guesses[strings.ToLower(string(letter))]}
	}

	// Données à passer au template HTML
	data := struct {
		Word        []struct{ Letter string; Guessed bool }
		Guesses     map[string]bool
		Message     string
		Hangman     string
		HangmanState int
	}{Word: wordWithState, Guesses: guesses, Message: message, HangmanState: hangmanState, Hangman: hangmanStates[hangmanState]}

	// Charger le template HTML
	tmpl, err := template.ParseFiles("./serv/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Afficher le template avec les données
	tmpl.Execute(w, data)
}

// Fonction principale du programme
func main() {
	// Sélectionner un mot à deviner
	wordToGuess = WordPicker(RandomNumber())
	if wordToGuess == "" {
		fmt.Println("Aucun mot trouvé. Arrêt du programme.")
		return
	}

	// Gestion des fichiers statiques
	fsServ := http.FileServer(http.Dir("serv"))
	http.Handle("/serv/", http.StripPrefix("/serv/", fsServ))

	fsImg := http.FileServer(http.Dir("img"))
	http.Handle("/img/", http.StripPrefix("/img/", fsImg))

	// Définir les gestionnaires d'URL
	http.HandleFunc("/", handler)          // Page principale du jeu
	http.HandleFunc("/restart", restartHandler)  // Redémarrer le jeu
	http.HandleFunc("/nextpage", nextPageHandler) // Passer à la page suivante (suite.html)

	// Nouvelle route pour suite.html
	http.HandleFunc("/suite", suiteHandler)

	// Démarrer le serveur HTTP
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erreur lors du démarrage du serveur:", err)
	}
}

// Gestionnaire d'URL pour la page principale du jeu
func handler(w http.ResponseWriter, r *http.Request) {
	allGuessed := true
	wordWithState := make([]struct {
		Letter  string
		Guessed bool
	}, len(wordToGuess))

	// Gérer les lettres devinées par l'utilisateur
	if r.Method == "POST" {
		r.ParseForm()
		guess := strings.ToLower(r.FormValue("guess"))
		if _, ok := guesses[guess]; !ok && !strings.Contains(wordToGuess, guess) {
			errorCount++
			hangmanState++
		}
		if hangmanState == 6 {
			hangmanState = 6
		}
		guesses[guess] = true
	}

	// Remplir le tableau wordWithState avec les lettres devinées et non devinées
	for i, letter := range wordToGuess {
		wordWithState[i] = struct {
			Letter  string
			Guessed bool
		}{Letter: string(letter), Guessed: guesses[strings.ToLower(string(letter))]}
	}

	// Vérifier si toutes les lettres ont été devinées ou si le nombre maximal d'erreurs a été atteint
	for _, letterState := range wordWithState {
		if !letterState.Guessed {
			allGuessed = false
			break
		}
	}

	// Afficher le message approprié et redémarrer le jeu si nécessaire
	if allGuessed || errorCount >= 6 {
		var message string
		if allGuessed {
			message = fmt.Sprintf("Félicitations ! Vous avez deviné le mot : %s", wordToGuess)
		} else {
			message = fmt.Sprintf("Désolé, vous avez dépassé le nombre maximal d'erreurs. Le mot était : %s", wordToGuess)
		}

		renderTemplate(w, message)

		resetGame()

		return
	}

	// Afficher le template de jeu
	renderTemplate(w, "")
}

// Fonction pour réinitialiser le jeu
func resetGame() {
	guesses = make(map[string]bool)
	wordToGuess = WordPicker(RandomNumber())
	errorCount = 0
	hangmanState = 0
}

// Gestionnaire d'URL pour redémarrer le jeu
func restartHandler(w http.ResponseWriter, r *http.Request) {
	resetGame()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Gestionnaire d'URL pour passer à la page suivante (suite.html)
func nextPageHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/suite", http.StatusSeeOther)
}

// Gestionnaire d'URL pour suite.html
func suiteHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./serv/suite1.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// Fonction pour sélectionner un mot aléatoire dans le fichier words.txt
func WordPicker(lineNumber int) string {
	file, err := os.Open("words.txt")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier words.txt:", err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0

	for scanner.Scan() {
		if currentLine == lineNumber {
			return scanner.Text()
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier words.txt:", err)
	}

	return ""
}

// Fonction pour générer un nombre aléatoire en fonction du nombre total de mots dans le fichier words.txt
func RandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(Counter())
}

// Fonction pour compter le nombre total de mots dans le fichier words.txt
func Counter() int {
	file, err := os.Open("words.txt")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier words.txt:", err)
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier words.txt:", err)
	}

	return count
}
