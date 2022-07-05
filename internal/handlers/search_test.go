package handlers

import (
	"fmt"
	"ghebant/lbc-api/internal/levenshtein"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
)

//func TestNormalizeKeywords(t *testing.T) {
//	source := NormalizeKeywords([]string{"Bien", "joué"})
//
//	assert.True(t, len(source) == 2)
//	assert.Equal(t, "bien", source[0])
//	assert.Equal(t, "joue", source[1])
//}
//

func TestLongestCommonSubstrings(t *testing.T) {
	needle := "zabcd"
	word := "mabzd"

	LongestCommonSubstrings(needle, word)
}

func TestAA(t *testing.T) {
	search := "avant"

	nb := levenshtein.Levenshtein(search, "s4 avant")
	log.Println(nb)
}

type Input struct {
	word string
	// nb changements
	score int
}

type BestMatch struct {
	brand  string
	model  string
	inputs []Input
	score  int
}

func Test2(t *testing.T) {
	Vehicles := map[string][]string{
		"Audi":    {"Cabriolet", "Q2", "Q3", "Q5", "Q7", "Q8", "R8", "Rs3", "Rs4", "Rs5", "Rs7", "S3", "S4", "S4 Avant", "S4 Cabriolet", "S5", "S7", "S8", "SQ5", "SQ7", "Tt", "Tts", "V8"},
		"BMW":     {"M3", "M4", "M5", "M535", "M6", "M635", "Serie 1", "Serie 2", "Serie 3", "Serie 4", "Serie 5", "Serie 6", "Serie 7", "Serie 8"},
		"Citroen": {"C1", "C15", "C2", "C25", "C25D", "C25E", "C25TD", "C3", "C3 Aircross", "C3 Picasso", "C4", "C4 Picasso", "C5", "C6", "C8", "Ds3", "Ds4", "Ds5"},
	}

	search := "rs4 avant"
	search = strings.ToLower(search)

	inputs := strings.Split(search, " ")

	//VehiclesScore := map[string]map[string]int{}

	VehiclesScore := make(map[string]map[string][]Input)

	best := BestMatch{"", "", []Input{}, 100}

exit:
	for brand, models := range Vehicles {
		for _, model := range models {

			// Search correspond exactly
			if model == search {
				best.score = 0
				best.model = model
				best.brand = brand

				continue exit
			}

			for _, input := range inputs {
				model = strings.ToLower(model)

				nb := levenshtein.Levenshtein(input, model)

				if VehiclesScore[brand] == nil {
					VehiclesScore[brand] = make(map[string][]Input)
				}

				scoredInput := Input{
					word:  input,
					score: nb,
				}

				VehiclesScore[brand][model] = append(VehiclesScore[brand][model], scoredInput)

				//VehiclesScore[brand][model] += nb
			}
		}
	}

	var winners []BestMatch

	if best.score != 0 {

		for brand, models := range VehiclesScore {
			log.Println(brand)
			for modelName, inputs := range models {

				// inputs == VehiclesScore[brand][modelName]

				globalScore := CalculateGlobalScore(inputs)

				if globalScore <= best.score {
					best.score = globalScore
					best.model = modelName
					best.brand = brand
					best.inputs = inputs
					winners = GetTop3(best, winners)
				}

				log.Println(inputs)
				log.Printf("  %s %d", modelName, globalScore)
			}
		}
	}

	winner := FindWinner(winners)
	log.Println("winners", winners)

	log.Println("winner", winner)
	//log.Printf("FOUND %s %s \n score: %d match: %s", best.brand, best.model, best.score, best.match)
}

func FindWinner(winners []BestMatch) BestMatch {
	bestInputScore := 100
	bestWinnerIndex := 0
	for i := range winners {
		for j := range winners[i].inputs {
			if winners[i].inputs[j].score < bestInputScore {
				bestInputScore = winners[i].inputs[j].score
				bestWinnerIndex = i
			}
		}
	}

	return winners[bestWinnerIndex]
}

func GetTop3(match BestMatch, bestMatches []BestMatch) []BestMatch {
	if len(bestMatches) < 3 || len(bestMatches) == 0 {
		bestMatches = append(bestMatches, match)
		return bestMatches
	}

	// Find worst score to replace it
	worstScoreIndex := 100
	worstScore := -1
	for i, bestMatch := range bestMatches {
		if bestMatch.score > worstScore {
			worstScore = bestMatch.score
			worstScoreIndex = i
		}
	}

	if match.score < bestMatches[worstScoreIndex].score {
		bestMatches[worstScoreIndex] = match
	}

	return bestMatches
}

func TestGetTop3(t *testing.T) {
	tests := []struct {
		match       BestMatch
		bestMatches []BestMatch
		expected    []BestMatch
	}{
		{BestMatch{score: 3}, []BestMatch{{score: 1}, {score: 4}, {score: 2}}, []BestMatch{{score: 1}, {score: 3}, {score: 2}}},
		{BestMatch{score: 10}, []BestMatch{{score: 1}, {score: 4}, {score: 2}}, []BestMatch{{score: 1}, {score: 4}, {score: 2}}},
		{BestMatch{score: 10}, []BestMatch{{score: 1}, {score: 4}}, []BestMatch{{score: 1}, {score: 4}, {score: 10}}},
		{BestMatch{score: 10}, []BestMatch{}, []BestMatch{{score: 10}}},
	}

	for i := range tests {
		tests[i].bestMatches = GetTop3(tests[i].match, tests[i].bestMatches)

		assert.Equal(t, tests[i].expected, tests[i].bestMatches)
	}
}

func CalculateGlobalScore(inputs []Input) int {
	score := 0
	for i := range inputs {
		score += inputs[i].score
	}

	return score
}

func sort() {
	var n = []int{1, 39, 2, 9, 7, 54, 11}

	var isDone = false

	for !isDone {
		isDone = true
		var i = 0
		for i < len(n)-1 {
			if n[i] > n[i+1] {
				n[i], n[i+1] = n[i+1], n[i]
				isDone = false
			}
			i++
		}
	}

	fmt.Println(n)
}
func Test(t *testing.T) {
	//str := "abcd"
	//needle := "zabcd"
	//needle2 := "mabzd"

	Vehicles := map[string][]string{
		"Audi":    {"Cabriolet", "Q2", "Q3", "Q5", "Q7", "Q8", "R8", "Rs3", "Rs4", "Rs5", "Rs7", "S3", "S4", "S4 Avant", "S4 Cabriolet", "S5", "S7", "S8", "SQ5", "SQ7", "Tt", "Tts", "V8"},
		"BMW":     {"M3", "M4", "M5", "M535", "M6", "M635", "Serie 1", "Serie 2", "Serie 3", "Serie 4", "Serie 5", "Serie 6", "Serie 7", "Serie 8"},
		"Citroen": {"C1", "C15", "C2", "C25", "C25D", "C25E", "C25TD", "C3", "C3 Aircross", "C3 Picasso", "C4", "C4 Picasso", "C5", "C6", "C8", "Ds3", "Ds4", "Ds5"},
	}

	VehiclesScore := map[string]map[string]int{}

	input := "Serie5"

	type bestMatch struct {
		match string
		brand string
		model int
	}

	type bestM struct {
		match string
		brand string
		model string
		score int
	}

	best := bestM{"", "", "", 0}

	for brand, models := range Vehicles {
		for _, model := range models {
			model = strings.ToLower(model)

			for n := 0; n <= len(input); n++ {

				for j := 0; j < len(input); j++ {
					if n+j <= len(input) {
						nb := pdr(model, input[j:j+n])

						VehiclesScore[brand][model] += 1
						//if nb > best.score {
						//	best.score = nb
						//	best.match = input[j : j+n]
						//	best.model = model
						//	best.brand = brand
						//}

						log.Println(model, input[j:j+n], nb)
					}
				}

			}
		}
	}

	for brand, _ := range VehiclesScore {
		for model, _ := range VehiclesScore {
			if VehiclesScore[brand][model] > best.score {
				best.score = VehiclesScore[brand][model]
				best.model = model
				best.brand = brand
			}
		}
	}

	log.Printf("FOUND %s %s \n score: %d match: %s", best.brand, best.model, best.score, best.match)
	//log.Printf("FOUND %s %s \n score: %d match: %s", best.brand, best.model, best.score, best.match)

	//for n := 0; n <= len(str); n++ {
	//
	//	for j := 0; j < len(str); j++ {
	//		if n+j <= len(str) {
	//			nb := pdr(needle2, str[j:j+n])
	//			log.Println(str[j:j+n], nb)
	//		}
	//	}
	//}
}

//
//func TestLcs2(t *testing.T) {
//	needle := "zabcd"
//	word := "mabzd"
//
//	log.Println(Lcs2(needle, word))
//}

//func TestPostAd2(t *testing.T) {
//	str1 := "OldSite:GeeksforGeeks.org"
//	str2 := "NewSite:GeeksQuiz.com"
//
//	res := pdr(str1, str2)
//	log.Println(res)
//}
