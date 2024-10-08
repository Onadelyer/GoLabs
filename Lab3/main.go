package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	inputFileName := "input.txt"
	outputFileName := "output.txt"

	data, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		fmt.Println("Помилка читання файлу:", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	var outputLines []string

	for _, line := range lines {
		words := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == ','
		})
		var newWords []string
		for _, word := range words {
			newWord := replaceDoubleLetters(word)
			fmt.Println("Слово з заміненими подвоєними літерами:", newWord)
			shuffledWord := shuffleString(newWord)
			fmt.Println("Слово з перемішаними літерами:", shuffledWord)
			newWords = append(newWords, shuffledWord)
		}
		outputLines = append(outputLines, strings.Join(newWords, "-"))
	}

	outputData := strings.Join(outputLines, "\n")

	err = ioutil.WriteFile(outputFileName, []byte(outputData), 0644)
	if err != nil {
		fmt.Println("Помилка запису файлу:", err)
		return
	}

	fmt.Println("Обробку завершено. Перевірте", outputFileName)
}

func replaceDoubleLetters(word string) string {
	runes := []rune(word)
	for i := 0; i < len(runes)-1; i++ {
		if runes[i] == runes[i+1] {
			runes[i] = '+'
			runes = append(runes[:i+1], runes[i+2:]...)
			i--
		}
	}
	return string(runes)
}

func shuffleString(input string) string {
	runes := []rune(input)
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}