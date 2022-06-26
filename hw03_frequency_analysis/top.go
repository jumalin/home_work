package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type scoreWordsMapping struct {
	score int
	words []string
}

var regex = regexp.MustCompile(`[\s\t\n.,!?'_:;\]\[(){}]`)

const ten = 10

func Top10(text string) (top10 []string) {
	if text == "" {
		return []string{}
	}

	scoreWordsSlice := scoreWordsSlice(text)
	for _, sw := range scoreWordsSlice {
		sort.Strings(sw.words)
		top10 = append(top10, sw.words[:min(len(sw.words), ten-len(top10))]...)

		if len(top10) == ten {
			break
		}
	}
	return
}

func scoreWordsSlice(text string) []scoreWordsMapping {
	words := regex.Split(text, -1)

	wordScoreMap := make(map[string]int, len(words))
	for _, word := range words {
		if word == "-" || word == "" {
			continue
		}
		wordScoreMap[strings.ToLower(word)]++
	}

	scoreWordsMap := make(map[int][]string, len(words))
	for word, score := range wordScoreMap {
		words := scoreWordsMap[score]
		scoreWordsMap[score] = append(words, word)
	}

	scores := make([]scoreWordsMapping, 0, len(words))
	for score, words := range scoreWordsMap {
		scores = append(scores, scoreWordsMapping{
			words: words,
			score: score,
		})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})
	return scores[:min(len(scores), ten)]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
