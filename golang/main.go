package main

import (
	"html/template"
	"net/http"
	"strings"
)

var (
	wordToGuess = "word.txt"
	guesses     = make(map[string]bool)
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		guess := strings.ToLower(r.FormValue("guess"))
		guesses[guess] = true
	}

	renderTemplate(w)
}

func renderTemplate(w http.ResponseWriter) {
	tmpl, err := template.New("index.html").Parse(`<!doctype html>
	<html>
	<head>
	  <meta charset="utf-8">
	  <meta name="keywords" content="Altrévis">
	  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	  <!-- CSS perso-->
	  <link rel="stylesheet" href="style.css" />
	  <title>Hangman web rétro</title>
	</head>
	<body>
	
	  <header>
	
		<nav class="navbar">
		  <div class="navbar-container container">
			<input type="checkbox" name="" id="" />
			<div class="hamburger-lines">
			  <span class="line line1"></span>
			  <span class="line line2"></span>
			  <span class="line line3"></span>
			</div>
			<ul class="menu-items">
			  <li><a href="index.html">Menu</a></li>
			  <li><a href="rule.html">Règle</a></li>
			</ul>
			<h2></h2>
		  </div>
		</nav>
	
		<video autoplay muted loop>
		  <source src="img/.mp4" type="video/mp4">
		</video>
	
		<div class="content">
		  <h1><a class="button" href="play.html">Jouer</a></h1>
		  <h3></h3>
		</div>
	
	  </header>
	
	</body>
	</html>`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wordWithState := make([]struct {
		Letter  string
		Guessed bool
	}, len(wordToGuess))

	for i, letter := range wordToGuess {
		wordWithState[i] = struct {
			Letter  string
			Guessed bool
		}{Letter: string(letter), Guessed: guesses[strings.ToLower(string(letter))]}
	}

	data := struct {
		Word    []struct{ Letter string; Guessed bool }
		Guesses map[string]bool
	}{Word: wordWithState, Guesses: guesses}

	tmpl.Execute(w, data)
}