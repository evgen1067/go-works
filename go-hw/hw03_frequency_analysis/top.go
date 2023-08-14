package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Word struct {
	value     string
	frequency int
}

func textToMap(s string) map[string]int {
	wordsMap := make(map[string]int)
	fields := strings.Fields(s)
	for _, word := range fields {
		lowerCaseWord := strings.ToLower(word)
		cleanString(&lowerCaseWord)
		_, ok := wordsMap[lowerCaseWord]
		if ok {
			wordsMap[lowerCaseWord]++
		} else {
			wordsMap[lowerCaseWord] = 1
		}
	}
	return wordsMap
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func cleanString(sPtr *string) {
	punctuation := []string{",", ".", ";", "?", "!", "'", `"`}
	for i := range punctuation {
		if strings.Contains(*sPtr, punctuation[i]) {
			*sPtr = strings.ReplaceAll(*sPtr, punctuation[i], "")
		}
	}
}

func sortingWords(words []Word) []string {
	sort.Slice(words, func(i, j int) bool {
		return words[i].frequency > words[j].frequency ||
			(words[i].frequency == words[j].frequency && words[i].value < words[j].value)
	})
	minimum := min(10, len(words))
	slice := make([]string, minimum)
	for i := range slice {
		slice[i] = words[i].value
	}
	return slice
}

func Top10(s string) []string {
	wordsMap := textToMap(s)
	var words []Word
	for value, frequency := range wordsMap {
		if value != "-" {
			words = append(words, Word{value, frequency})
		}
	}
	slice := sortingWords(words)
	return slice
}
