package matching

import (
	"ghebant/lbc-api/internal/helpers"
	"ghebant/lbc-api/internal/levenshtein"
	"ghebant/lbc-api/models"
	"math"
	"strings"
)

func FindBestMatch(search string, Vehicles *map[string][]string) models.Match {
	search = helpers.NormalizeString(search)
	searchKeywords := strings.Split(search, " ")
	currentBestMatch := models.Match{KeywordWithScores: []models.KeywordWithScore{}, Distance: 100}

	// TODO make func ComputeGlobalScores()
	for carBrand, carModels := range *Vehicles {
		for _, carModel := range carModels {
			carModel = helpers.NormalizeString(carModel)

			// If search correspond exactly
			if carModel == search {
				currentBestMatch.Distance = 0
				currentBestMatch.Model = carModel
				currentBestMatch.Brand = carBrand

				// Finner found
				return currentBestMatch
			}

			keywordsWithScores := ComputeKeywordsScores(carModel, searchKeywords)

			globalDistance, averageMatchingPercent := ComputeGlobalScores(keywordsWithScores)

			// Replace old best match by the new one
			if averageMatchingPercent >= currentBestMatch.AverageMatchingPercent {
				if averageMatchingPercent == currentBestMatch.AverageMatchingPercent {
					// TODO
				}

				currentBestMatch.Distance = globalDistance
				currentBestMatch.Model = carModel
				currentBestMatch.Brand = carBrand
				currentBestMatch.KeywordWithScores = keywordsWithScores
				currentBestMatch.AverageMatchingPercent = averageMatchingPercent
			}
		}
	}

	return currentBestMatch
}

// ComputeKeywordsScores computes and attributes scores to each keyword
func ComputeKeywordsScores(carModel string, keywords []string) []models.KeywordWithScore {
	var keywordsWithScores []models.KeywordWithScore

	for i := range keywords {
		// For each keyword compute Distance between the keyword and the car Model
		distance := levenshtein.Levenshtein(keywords[i], carModel)

		// Compute how much percentage of the search keyword is present in the car Model
		matchingPercentage := ComputeMatchingPercentage(len(carModel), distance)

		keywordWithScore := models.KeywordWithScore{
			Keyword:            keywords[i],
			DistanceFromModel:  distance,
			MatchingPercentage: matchingPercentage,
		}

		keywordsWithScores = append(keywordsWithScores, keywordWithScore)
	}

	return keywordsWithScores
}

// ComputeGlobalScores returns the global Distance and the global average matching percentage from all the keywords
func ComputeGlobalScores(keywords []models.KeywordWithScore) (int, float64) {
	globalDistance := 0
	averageMatchingPercentage := 0.0

	for i := range keywords {
		globalDistance += keywords[i].DistanceFromModel

		if keywords[i].MatchingPercentage > 0 {
			averageMatchingPercentage += keywords[i].MatchingPercentage
		}
	}

	if averageMatchingPercentage > 0 {
		averageMatchingPercentage = averageMatchingPercentage / float64(len(keywords))
		averageMatchingPercentage = math.Round(averageMatchingPercentage*100) / 100
	}

	return globalDistance, averageMatchingPercentage
}

// ComputeMatchingPercentage computes the percentage of a Keyword A matching another Keyword B based on Keyword B length and Keyword's A Distance from Keyword B
func ComputeMatchingPercentage(modelLength, distanceFromWord int) float64 {
	if modelLength > 0 {
		matchPercent := 100 - (float64(distanceFromWord)*100)/float64(modelLength)
		return math.Round(matchPercent*100) / 100
	}

	return 0
}
