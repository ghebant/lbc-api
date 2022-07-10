package matching

import (
	"ghebant/lbc-api/internal/constants"
	"ghebant/lbc-api/internal/helpers"
	"ghebant/lbc-api/internal/levenshtein"
	"ghebant/lbc-api/models"
	"math"
	"strings"
)

// FindBestMatch returns the best match found for 'search' in all Vehicles
func FindBestMatch(search string, Vehicles *map[string][]string) models.Match {
	search = helpers.NormalizeString(search)
	currentBestMatch := models.Match{KeywordWithScores: []models.KeywordWithScore{}}

	for carBrand, carModels := range *Vehicles {
		for _, carModel := range carModels {
			carModel = helpers.NormalizeString(carModel)

			match := ComputeMatch(search, carModel)

			if match.AverageMatchingPercent == 1000 {
				match.Brand = carBrand
				return match
			}

			if match.AverageMatchingPercent > currentBestMatch.AverageMatchingPercent {
				currentBestMatch = match
				currentBestMatch.Brand = carBrand
			}
		}
	}

	return currentBestMatch
}

// ComputeMatch Returns an average percentage for 'search' matching 'carModel'
func ComputeMatch(search, carModel string) models.Match {
	match := models.Match{}
	searchKeywords := strings.Split(search, " ")

	// If search correspond exactly to the model return the winner
	if carModel == search {
		match.Model = carModel
		match.AverageMatchingPercent = 1000

		return match
	}

	// Compute scores for every search keyword
	keywordsWithScores := ComputeKeywordsScores(carModel, searchKeywords)

	// Compute global
	averageMatchingPercent := ComputeGlobalScores(keywordsWithScores)

	// Replace old best match by the new one
	if averageMatchingPercent >= match.AverageMatchingPercent {
		match.Model = carModel
		match.KeywordWithScores = keywordsWithScores
		match.AverageMatchingPercent = averageMatchingPercent
	}

	return match
}

// ComputeGlobalScores returns the global average matching percentage from all the keywords
func ComputeGlobalScores(keywords []models.KeywordWithScore) float64 {
	averageMatchingPercentage := 0.0

	for i := range keywords {
		if keywords[i].MatchingPercentage > 0 {
			averageMatchingPercentage += keywords[i].MatchingPercentage
		}
	}

	if averageMatchingPercentage > 0 {
		averageMatchingPercentage = averageMatchingPercentage / float64(len(keywords))
		averageMatchingPercentage = math.Round(averageMatchingPercentage*100) / 100
	}

	return averageMatchingPercentage
}

// ComputeKeywordsScores computes and attributes scores to each search keyword, the final MatchingPercentage can go above 100% du to the weight added
func ComputeKeywordsScores(carModel string, keywords []string) []models.KeywordWithScore {
	var keywordsWithScores []models.KeywordWithScore

	carModelKeywords := strings.Split(carModel, " ")

	for i := range keywords {
		bestAverage := 0.0

		// Compute the distance between a search keyword and the whole model
		distanceFromWholeModel := levenshtein.Levenshtein(keywords[i], carModel)
		mpFromModel := ComputeMatchingPercentage(len(carModel), distanceFromWholeModel)

		for j := range carModelKeywords {

			// Compute the distance between a search keyword and a keyword from the car model
			distance := levenshtein.Levenshtein(keywords[i], carModelKeywords[j])

			// Compute how much percentage of the search keyword matches the car model keyword based on there distance
			matchingPercentage := ComputeMatchingPercentage(len(carModelKeywords[j]), distance)

			// Returns a percentage based on how many characters in a row they have in common since index 0
			percentageCharactersMatching := ComputeCharactersMatching(keywords[i], carModelKeywords[j])

			// Calculate the weight of a keyword based on its length
			weight := CalculateKeywordWeight(keywords[i])

			globalAverage := ((matchingPercentage * weight) + (percentageCharactersMatching * weight) + mpFromModel) / 3

			if globalAverage > bestAverage {
				bestAverage = globalAverage
			}

		}

		keywordWithScore := models.KeywordWithScore{
			Keyword:            keywords[i],
			MatchingPercentage: bestAverage,
		}

		keywordsWithScores = append(keywordsWithScores, keywordWithScore)
	}

	return keywordsWithScores
}

// ComputeMatchingPercentage computes the percentage of a keyword matching another one based on its length, and it's distance from it
func ComputeMatchingPercentage(modelLength, distanceFromModel int) float64 {
	if distanceFromModel > modelLength {
		return 0
	}

	if modelLength > 0 {
		matchPercent := 100 - (float64(distanceFromModel)*100)/float64(modelLength)
		return math.Round(matchPercent*100) / 100
	}

	return 0
}

// ComputeCharactersMatching Compute how much percent of a string matches another one base on how many characters in a row they have in common since index 0
func ComputeCharactersMatching(str1, str2 string) float64 {
	matchingCharacters := 0.0
	percentage := 0.0
	shortestLength := len(str1)
	longestLength := len(str2)

	if len(str1) > len(str2) {
		shortestLength = len(str2)
		longestLength = len(str1)
	}

	for i := 0; i < shortestLength; i++ {
		if str1[i] != str2[i] {
			break
		}

		matchingCharacters += 1
	}

	if matchingCharacters > 0 && longestLength > 0 {
		percentage = matchingCharacters * 100 / float64(longestLength)
	}

	return percentage
}

// CalculateKeywordWeight Calculates the weight of a word based on its length so that the longer the word is the more weight it has
func CalculateKeywordWeight(keyword string) float64 {
	weight := 1.0

	if len(keyword) > 1 {
		weight = weight + float64(len(keyword))*constants.KeywordWeight
	}

	return weight
}
