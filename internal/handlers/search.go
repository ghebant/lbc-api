package handlers

import (
	"ghebant/lbc-api/internal/levenshtein"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"math"
	"net/http"
	"strings"
	"unicode"
)

// TODO
// Lowercase everything

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		// TODO add error management
		panic(e)
	}
	return output
}

func NormalizeString(str string) string {
	// lower case
	str = strings.ToLower(str)
	// remove accents
	str = RemoveAccents(str)

	return str
}

type SearchKeyword struct {
	Keyword string
	// nb changements
	DistanceFromModel  int
	MatchingPercentage float64
}

type BestMatch struct {
	brand                  string
	model                  string
	inputs                 []SearchKeyword
	distance               int
	averageMatchingPercent float64
}

func SearchCar(search string) BestMatch {
	Vehicles := map[string][]string{
		"Audi":    {"Cabriolet", "Q2", "Q3", "Q5", "Q7", "Q8", "R8", "Rs3", "Rs4", "Rs5", "Rs7", "S3", "S4", "S4 Avant", "S4 Cabriolet", "S5", "S7", "S8", "SQ5", "SQ7", "Tt", "Tts", "V8"},
		"BMW":     {"M3", "M4", "M5", "M535", "M6", "M635", "Serie 1", "Serie 2", "Serie 3", "Serie 4", "Serie 5", "Serie 6", "Serie 7", "Serie 8"},
		"Citroen": {"C1", "C15", "C2", "C25", "C25D", "C25E", "C25TD", "C3", "C3 Aircross", "C3 Picasso", "C4", "C4 Picasso", "C5", "C6", "C8", "Ds3", "Ds4", "Ds5"},
	}

	search = NormalizeString(search)

	searchKeywords := strings.Split(search, " ")

	VehiclesScore := make(map[string]map[string][]SearchKeyword)

	currentBestMatch := BestMatch{"", "", []SearchKeyword{}, 100, 0}

	//var bestMatches []BestMatch

	// TODO make func ComputeScores()
	for carBrand, carModels := range Vehicles {
		for _, carModel := range carModels {
			carModel = NormalizeString(carModel)

			// If search correspond exactly
			if carModel == search {
				currentBestMatch.distance = 0
				currentBestMatch.model = carModel
				currentBestMatch.brand = carBrand

				// Finner found
				return currentBestMatch
			}

			var keywordsWithScores []SearchKeyword

			for i := range searchKeywords {
				// For each keyword compute distance between the keyword and the car model
				distance := levenshtein.Levenshtein(searchKeywords[i], carModel)

				if VehiclesScore[carBrand] == nil {
					VehiclesScore[carBrand] = make(map[string][]SearchKeyword)
				}

				// Compute how much percentage of the search keyword is present in the car model
				matchingPercentage := ComputeMatchingPercentage(len(carModel), distance)

				keywordWithScore := SearchKeyword{
					Keyword:            searchKeywords[i],
					DistanceFromModel:  distance,
					MatchingPercentage: matchingPercentage,
				}

				//VehiclesScore[carBrand][carModel] = append(VehiclesScore[carBrand][carModel], scoredInput)
				keywordsWithScores = append(keywordsWithScores, keywordWithScore)
			}

			globalDistance, averageMatchingPercent := ComputeScores(keywordsWithScores)

			if averageMatchingPercent >= currentBestMatch.averageMatchingPercent {
				if averageMatchingPercent == currentBestMatch.averageMatchingPercent {
					// TODO
				}

				currentBestMatch.distance = globalDistance
				currentBestMatch.model = carModel
				currentBestMatch.brand = carBrand
				currentBestMatch.inputs = keywordsWithScores
				currentBestMatch.averageMatchingPercent = averageMatchingPercent
				// Replace the worst match in the top 3 list by the new best match
				//bestMatches = Top3BestMatches(currentBestMatch, bestMatches)
			}
		}
	}

	winner := currentBestMatch
	//log.Println("bestMatches", bestMatches)
	//log.Println("winner", winner)

	return winner
}

// ComputeScores returns the global distance and the global average matching percentage of the keywords
func ComputeScores(keywords []SearchKeyword) (int, float64) {
	globalDistance := 0
	averageMatchingPercentage := 0.0

	for i := range keywords {
		globalDistance += keywords[i].DistanceFromModel

		if keywords[i].MatchingPercentage > 0 {
			averageMatchingPercentage += keywords[i].MatchingPercentage
		}
	}

	averageMatchingPercentage = averageMatchingPercentage / float64(len(keywords))

	return globalDistance, averageMatchingPercentage
}

func Top3BestMatches(match BestMatch, bestMatches []BestMatch) []BestMatch {
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

// Calculates the percentage of a Keyword A matching another Keyword B based on Keyword B length and Keyword's A distance from Keyword B
func ComputeMatchingPercentage(modelLength, distanceFromWord int) float64 {
	matchPercent := 100 - (float64(distanceFromWord)*100)/float64(modelLength)
	return math.Round(matchPercent*100) / 100
}

func Search(c *gin.Context) {
	queryParam := c.Param("keywords")

	if len(queryParam) <= 0 {
		log.Println("invalid keywords")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid keywords"})
		return
	}

	//keywords := strings.Split(queryParam, ",")

	log.Println("distance entre maison et macon:", levenshtein.Levenshtein("maison", "macon"))

	// TODO
	// Remove caps, spaces and accent

	c.JSON(http.StatusOK, nil)
}
