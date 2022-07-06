package handlers

import (
	"ghebant/lbc-api/internal/levenshtein"
	"github.com/stretchr/testify/assert"
	"log"
	"math"
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
	score        int
	matchPercent float64
}

type BestMatch struct {
	brand                  string
	model                  string
	inputs                 []Input
	score                  int
	averageMatchingPercent float64
}

func TestSearchCar(t *testing.T) {
	tests := []struct {
		search        string
		expectedBrand string
		expectedModel string
	}{
		{"rs4 avant", "Audi", "rs4"},
		{"Gran Turismo Serie5", "BMW", "serie 5"},
		{"ds 3 crossback", "Citroen", "ds3"},
		{"crossback ds 3", "Citroen", "ds3"},
		{"test avant", "Audi", "s4 avant"},
		{"cabr", "Audi", "cabriolet"},
		{"s4 avant", "Audi", "s4 avant"},
		// TODO Donner point si sub present au debut
		//{"s4 c", "Audi", "s4 cabriolet"},
	}

	for i := range tests {
		found := SearchCar(tests[i].search)

		assert.Equal(t, tests[i].expectedBrand, found.brand)
		assert.Equal(t, tests[i].expectedModel, found.model)
	}
}

func SearchCar(search string) BestMatch {
	Vehicles := map[string][]string{
		"Audi":    {"Cabriolet", "Q2", "Q3", "Q5", "Q7", "Q8", "R8", "Rs3", "Rs4", "Rs5", "Rs7", "S3", "S4", "S4 Avant", "S4 Cabriolet", "S5", "S7", "S8", "SQ5", "SQ7", "Tt", "Tts", "V8"},
		"BMW":     {"M3", "M4", "M5", "M535", "M6", "M635", "Serie 1", "Serie 2", "Serie 3", "Serie 4", "Serie 5", "Serie 6", "Serie 7", "Serie 8"},
		"Citroen": {"C1", "C15", "C2", "C25", "C25D", "C25E", "C25TD", "C3", "C3 Aircross", "C3 Picasso", "C4", "C4 Picasso", "C5", "C6", "C8", "Ds3", "Ds4", "Ds5"},
	}

	search = strings.ToLower(search)

	inputs := strings.Split(search, " ")

	//VehiclesScore := map[string]map[string]int{}

	VehiclesScore := make(map[string]map[string][]Input)

	best := BestMatch{"", "", []Input{}, 100, 0}

	for brand, models := range Vehicles {
		for _, model := range models {
			model = strings.ToLower(model)

			// Search correspond exactly
			if model == search {
				best.score = 0
				best.model = model
				best.brand = brand

				// Finner found
				return best
			}

			for _, input := range inputs {
				distance := levenshtein.Levenshtein(input, model)

				if VehiclesScore[brand] == nil {
					VehiclesScore[brand] = make(map[string][]Input)
				}

				matchPercent := CalculateMatchPercent(len(model), distance)

				scoredInput := Input{
					word:         input,
					score:        distance,
					matchPercent: matchPercent,
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

				globalScore, averageMatchingPercent := CalculateScores(inputs)

				if averageMatchingPercent >= best.averageMatchingPercent {
					//if globalScore <= best.score {
					best.score = globalScore
					best.model = modelName
					best.brand = brand
					best.inputs = inputs
					best.averageMatchingPercent = averageMatchingPercent
					//winners = GetTop3(best, winners)
					winners = GetTop3FromMatchingPercent(best, winners)
				}

				log.Println(inputs)
				log.Printf("  %s distance: %d percent matching: %f", modelName, globalScore, averageMatchingPercent)
			}
		}
	}

	//winner := FindWinner(winners)
	winner := FindWinner2(winners)

	//log.Println("winners", winners)
	//log.Println("winner", winner)

	return winner
}

// Calculates the percentage of a word A matching another word B based on word B length and word's A distance from word B
func CalculateMatchPercent(modelLength, distanceFromWord int) float64 {
	matchPercent := 100 - (float64(distanceFromWord)*100)/float64(modelLength)
	return math.Round(matchPercent*100) / 100
}

func TestCalculateMatchPercent(t *testing.T) {
	tests := []struct {
		modelLength int
		distance    int
		expected    float64
	}{
		{2, 1, 50},
		{3, 2, 33.33},
		{3, 0, 100},
		{3, 8, -166.67},
	}

	for i := range tests {
		assert.Equal(t, tests[i].expected, CalculateMatchPercent(tests[i].modelLength, tests[i].distance))
	}
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

func FindWinner2(winners []BestMatch) BestMatch {
	bestWinnerIndex := 0
	bestAverageMatching := 0.0

	for i := range winners {
		if winners[i].averageMatchingPercent > bestAverageMatching {
			bestAverageMatching = winners[i].averageMatchingPercent
			bestWinnerIndex = i
		}

		// If they have the same average matching
		if winners[i].averageMatchingPercent == bestAverageMatching {
			// take the one with the input with less changes
			// TODO
		}
	}

	return winners[bestWinnerIndex]
}

func TestGetTop3FromMatchingPercent(t *testing.T) {
	tests := []struct {
		match       BestMatch
		bestMatches []BestMatch
		expected    []BestMatch
	}{
		{
			BestMatch{averageMatchingPercent: 10},
			[]BestMatch{{averageMatchingPercent: 1}, {averageMatchingPercent: 4}, {averageMatchingPercent: 2}},
			[]BestMatch{{averageMatchingPercent: 10}, {averageMatchingPercent: 4}, {averageMatchingPercent: 2}},
		},
		{
			BestMatch{averageMatchingPercent: 38},
			[]BestMatch{{averageMatchingPercent: 40}, {averageMatchingPercent: 56}, {averageMatchingPercent: 89}},
			[]BestMatch{{averageMatchingPercent: 40}, {averageMatchingPercent: 56}, {averageMatchingPercent: 89}},
		},
		{
			BestMatch{averageMatchingPercent: 38},
			[]BestMatch{{averageMatchingPercent: 20}, {averageMatchingPercent: 56}, {averageMatchingPercent: 89}},
			[]BestMatch{{averageMatchingPercent: 38}, {averageMatchingPercent: 56}, {averageMatchingPercent: 89}},
		},
		{
			BestMatch{averageMatchingPercent: 38},
			[]BestMatch{{averageMatchingPercent: 20}, {averageMatchingPercent: 19}, {averageMatchingPercent: 18}},
			[]BestMatch{{averageMatchingPercent: 20}, {averageMatchingPercent: 19}, {averageMatchingPercent: 38}},
		},
	}

	for i := range tests {
		tests[i].bestMatches = GetTop3FromMatchingPercent(tests[i].match, tests[i].bestMatches)

		assert.Equal(t, tests[i].expected, tests[i].bestMatches)
	}
}

func GetTop3FromMatchingPercent(match BestMatch, bestMatches []BestMatch) []BestMatch {
	if len(bestMatches) < 3 || len(bestMatches) == 0 {
		bestMatches = append(bestMatches, match)
		return bestMatches
	}

	// TODO JE RECALCUL QUI EST LE MAILLEUR A CHAQUE FOIS -> PAS BON
	// TODO CALCULER ME MEILLEUR UNE FOIS ET COMPARER match AVEC
	// Find which input has the worst matching percent and replace it by the current match if it's higher
	worstMatchIndex := 0
	worstMatchingPercent := 1000.0
	for i := range bestMatches {
		if bestMatches[i].averageMatchingPercent < worstMatchingPercent {
			worstMatchIndex = i
			worstMatchingPercent = bestMatches[i].averageMatchingPercent
		}
	}

	if match.averageMatchingPercent > worstMatchingPercent {
		bestMatches[worstMatchIndex] = match
	}

	return bestMatches
}

func GetTop3(match BestMatch, bestMatches []BestMatch) []BestMatch {
	if len(bestMatches) < 3 || len(bestMatches) == 0 {
		bestMatches = append(bestMatches, match)
		return bestMatches
	}

	// TODO JE RECALCUL QUI EST LE MAILLEUR A CHAQUE FOIS -> PAS BON
	// TODO CALCULER ME MEILLEUR UNE FOIS ET COMPARER match AVEC
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

// return global score && global average matching percent of the inputs
func CalculateScores(inputs []Input) (int, float64) {
	score := 0
	averageMatchingPercent := 0.0

	for i := range inputs {
		score += inputs[i].score

		if inputs[i].matchPercent > 0 {
			averageMatchingPercent += inputs[i].matchPercent
		}
	}

	averageMatchingPercent = averageMatchingPercent / float64(len(inputs))

	return score, averageMatchingPercent
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
