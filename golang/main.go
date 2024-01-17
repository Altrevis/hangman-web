package main

import (
	// Importation des packages nécessaires
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"bufio"
)

// Variables globales
var (
	guesses      = make(map[string]bool) // Carte pour stocker les lettres devinées
	wordToGuess  string                   // Mot à deviner
	errorCount   int                      // Compteur d'erreurs
	hangmanState int                      // État actuel du pendu
)

// Tableau des différentes étapes visuelles du pendu
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


// Fonction pour rendre le template HTML avec les données fournies
func renderTemplate(w http.ResponseWriter, message string) {
	// Prépare les données à utiliser dans le template
	wordWithState := make([]struct {
		Letter  string
		Guessed bool
	}, len(wordToGuess))

	// Crée une structure de données pour le template
	data := struct {
		Word        []struct{ Letter string; Guessed bool }
		Guesses     map[string]bool
		Message     string
		Hangman     string
		HangmanState int
		HangmanStates []string
	}{
		Word:          wordWithState,
		Guesses:       guesses,
		Message:       message,
		Hangman:       hangmanStates[hangmanState],
		HangmanState:  hangmanState,
		HangmanStates: hangmanStates,
	}

	// Charge le template HTML depuis le fichier
	tmpl, err := template.ParseFiles("./serv/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Exécute le template avec les données et écrit le résultat dans la réponse HTTP
	tmpl.Execute(w, data)
}

// Fonction principale
func main() {
	// Initialisation du mot à deviner
	wordToGuess = WordPicker(RandomNumber())
	if wordToGuess == "" {
		fmt.Println("Aucun mot trouvé. Arrêt du programme.")
		return
	}

	// Configuration des gestionnaires d'URL
	fsServ := http.FileServer(http.Dir("serv"))
	http.Handle("/serv/", http.StripPrefix("/serv/", fsServ))

	fsImg := http.FileServer(http.Dir("img"))
	http.Handle("/img/", http.StripPrefix("/img/", fsImg))

	http.HandleFunc("/", handler)
	http.HandleFunc("/restart", restartHandler)
	http.HandleFunc("/nextpage", nextPageHandler)
	http.HandleFunc("/hangman", hangmanHandler)

	// Ajout d'une nouvelle route pour suite.html
	http.HandleFunc("/suite", suiteHandler)

	// Démarrage du serveur HTTP sur le port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erreur lors du démarrage du serveur:", err)
	}
}

// Fonction pour gérer les requêtes HTTP sur la page principale du jeu
func handler(w http.ResponseWriter, r *http.Request) {
	allGuessed := true
	wordWithState := make([]struct {
		Letter  string
		Guessed bool
	}, len(wordToGuess))

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

	for i, letter := range wordToGuess {
		wordWithState[i] = struct {
			Letter  string
			Guessed bool
		}{Letter: string(letter), Guessed: guesses[strings.ToLower(string(letter))]}
	}

	for _, letterState := range wordWithState {
		if !letterState.Guessed {
			allGuessed = false
			break
		}
	}

	if allGuessed || errorCount >= 6 {
		var message string
		if allGuessed {
			message = fmt.Sprintf("Félicitations ! Vous avez deviné le mot : %s", wordToGuess)
			// Redirigez vers suite(good).html sur victoire
			http.Redirect(w, r, "/suite/good", http.StatusSeeOther)
		} else {
			message = fmt.Sprintf("Désolé, vous avez dépassé le nombre maximal d'erreurs. Le mot était : %s", wordToGuess)
			// Redirigez vers suite(dead).html sur défaite
			http.Redirect(w, r, "/suite/dead", http.StatusSeeOther)
		}

		renderTemplate(w, message)

		resetGame()

		return
	}

	renderTemplate(w, "")
}

// Fonction pour réinitialiser l'état du jeu
func resetGame() {
	guesses = make(map[string]bool)
	wordToGuess = WordPicker(RandomNumber())
	errorCount = 0
	hangmanState = 0
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	resetGame()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func nextPageHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/suite", http.StatusSeeOther)
}

func suiteHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./serv/suite.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func hangmanHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("./serv/hangman.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, nil)
}

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

func RandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(Counter())
}

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
