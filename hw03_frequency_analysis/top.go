package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCounter struct {
	word  string
	count int
}

func Top10(source string) []string {
	wordSlice := strings.Fields(source)
	wordCounterMap := make(map[string]int)

	for _, val := range wordSlice {
		_, ok := wordCounterMap[val]
		if ok {
			wordCounterMap[val]++
		} else {
			wordCounterMap[val] = 1
		}
	}

	wordCounterList := make([]WordCounter, 0, len(wordCounterMap))

	for key, value := range wordCounterMap {
		wordCounterList = append(wordCounterList, WordCounter{key, value})
	}

	sort.Slice(wordCounterList, func(i, j int) bool {
		return wordCounterList[i].word < wordCounterList[j].word
	})

	sort.SliceStable(wordCounterList, func(i, j int) bool {
		return wordCounterList[i].count > wordCounterList[j].count
	})

	topCount := 10
	topSlice := []string{}

	for i, value := range wordCounterList {
		if i == topCount {
			break
		}
		topSlice = append(topSlice, value.word)
	}

	return topSlice
}
