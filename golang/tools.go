package main

import (
	"os"
	"math/rand"
	"time"
	"bufio"
	"fmt"
)

func Clear() { // Efface l'affichage du terminal
	os.Stdout.WriteString("\x1b[3;J\x1b[H\x1b[2J")
}

func Exit() { // Quitte le programme
	os.Exit(0)
}

func Counter() int { // Renvoie le nombre de mots dans le fichier txt
	file, _ := os.Open("words.txt") 
	scanner := bufio.NewScanner(file)
	var count int
	for scanner.Scan() { // Tant que le retour a la ligne est possible
		count++
	}
	return count
}

func RandomNumber() int { // Génère un nombre entre 0 et le nombre de mots dans le fichier txt
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(Counter())
	return randomNumber
}

func WordPicker(line_number_of_word int) string { // Fonction qui recherche le mot par rapport au nombre aléatoire généré par RandomNumber()
	word := ""
	file, _ := os.Open("words.txt") // Prend en valeur le fichier texte
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() { // Boucle qui va atteindre le mot qui correspond a la ligne du random_number
		if i == line_number_of_word {
			word = scanner.Text()
		}
		i++
	}
	return word
}

func Hidden(word string) { // Fonction qui compte le nombre de caractères dans le mot et qui les remplace par des underscore
	for i := 0; i < len(word); i++ {
		hidden_word = append(hidden_word , "_")
	}
	fmt.Println(hidden_word)
}