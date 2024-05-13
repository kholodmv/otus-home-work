package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

func Top10(text string) []string {
	words := strings.Fields(text)

	wordCountFreq := make(map[string]int)

	for _, word := range words {
		wordCountFreq[word]++
	}

	wordFreqPair := make([][2]string, 0, len(wordCountFreq))

	for word, freq := range wordCountFreq {
		wordFreqPair = append(wordFreqPair, [2]string{word, fmt.Sprintf("%d", freq)})
	}

	sort.Slice(wordFreqPair, func(i, j int) bool {
		if wordFreqPair[i][1] == wordFreqPair[j][1] {
			return wordFreqPair[i][0] < wordFreqPair[j][0]
		}
		return wordFreqPair[i][1] > wordFreqPair[j][1]
	})

	topWords := make([]string, 0, 10)

	for i := 0; i < len(wordFreqPair) && len(topWords) < 10; i++ {
		topWords = append(topWords, wordFreqPair[i][0])
	}

	return topWords
}
