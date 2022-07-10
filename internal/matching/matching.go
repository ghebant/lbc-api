package matching

import (
	"ghebant/lbc-api/internal/helpers"
	"ghebant/lbc-api/internal/levenshtein"
	"ghebant/lbc-api/models"
	"log"
	"math"
	"strings"
)

func FindBest(search, carModel string) models.Match {
	match := models.Match{}
	searchKeywords := strings.Split(search, " ")
	//carModelKeywords := strings.Split(carModel, " ")

	// If search correspond exactly
	if carModel == search {
		match.Distance = 0
		match.Model = carModel
		match.AverageMatchingPercent = 1000

		// Finner found
		return match
	}

	//keywordsWithScores := ComputeKeywordsScores(carModel, searchKeywords)
	keywordsWithScores := ComputeKeywordsScores2SPLIT(carModel, searchKeywords)

	globalDistance, averageMatchingPercent := ComputeGlobalScores(keywordsWithScores)

	//if averageMatchingPercent > 0 {
	//	// Divide average by nb keyword missing
	//	nbWordDifference := len(searchKeywords) - len(carModelKeywords)
	//	if nbWordDifference < 0 {
	//		nbWordDifference *= -1
	//	}
	//
	//	nbWordDifference += 1
	//
	//	nn := averageMatchingPercent / float64(nbWordDifference)
	//
	//	//newAverage = newAverage / float64(nbWordDifference)
	//	log.Printf("average %f nbDiff: %d new: %f", averageMatchingPercent, nbWordDifference, nn)
	//	averageMatchingPercent = nn
	//}

	//newAverage := (averageMatchingPercent + (percentMatchingBeginning * 1.4)) / 2

	//log.Printf("search: %s model: %s distance: %d average: %f", search, carModel, globalDistance, averageMatchingPercent)
	//log.Println(keywordsWithScores)

	// Replace old best match by the new one
	if averageMatchingPercent >= match.AverageMatchingPercent {
		if averageMatchingPercent == match.AverageMatchingPercent {
			// TODO
		}

		match.Distance = globalDistance
		match.Model = carModel
		match.KeywordWithScores = keywordsWithScores
		match.AverageMatchingPercent = averageMatchingPercent
	}

	return match
}

func FindBestMatch(search string, Vehicles *map[string][]string) models.Match {
	var top3 []models.Match

	search = helpers.NormalizeString(search)
	//searchKeywords := strings.Split(search, " ")
	currentBestMatch := models.Match{KeywordWithScores: []models.KeywordWithScore{}, Distance: 100}

	// TODO make func ComputeGlobalScores()
	for carBrand, carModels := range *Vehicles {
		for _, carModel := range carModels {
			carModel = helpers.NormalizeString(carModel)

			match := FindBest(search, carModel)
			top3 = GetTop3FromMatchingPercent(match, top3)

			if match.AverageMatchingPercent == 1000 {
				match.Brand = carBrand
				return match
			}

			if match.AverageMatchingPercent >= currentBestMatch.AverageMatchingPercent {
				currentBestMatch = match
				currentBestMatch.Brand = carBrand
			}

			//// If search correspond exactly
			//if carModel == search {
			//	currentBestMatch.Distance = 0
			//	currentBestMatch.Model = carModel
			//	currentBestMatch.Brand = carBrand
			//
			//	// Finner found
			//	return currentBestMatch
			//}
			//
			////keywordsWithScores := ComputeKeywordsScores(carModel, searchKeywords)
			//keywordsWithScores := ComputeKeywordsScores2SPLIT(carModel, searchKeywords)
			//
			//globalDistance, averageMatchingPercent := ComputeGlobalScores(keywordsWithScores)
			//
			////newAverage := (averageMatchingPercent + (percentMatchingBeginning * 1.4)) / 2
			//
			////log.Printf("search: %s model: %s distance: %d average: %f", search, carModel, globalDistance, averageMatchingPercent)
			////log.Println(keywordsWithScores)
			//
			//// Replace old best match by the new one
			//if averageMatchingPercent >= currentBestMatch.AverageMatchingPercent {
			//	if averageMatchingPercent == currentBestMatch.AverageMatchingPercent {
			//		// TODO
			//	}
			//
			//	currentBestMatch.Distance = globalDistance
			//	currentBestMatch.Model = carModel
			//	currentBestMatch.Brand = carBrand
			//	currentBestMatch.KeywordWithScores = keywordsWithScores
			//	currentBestMatch.AverageMatchingPercent = averageMatchingPercent
			//}
		}
	}

	log.Println("winner", currentBestMatch.Model, "percentage", currentBestMatch.AverageMatchingPercent)
	log.Println("top3")
	for i := range top3 {
		log.Println(top3[i].Model, "percentage", top3[i].AverageMatchingPercent)
	}

	return currentBestMatch
}

func GetTop3FromMatchingPercent(match models.Match, bestMatches []models.Match) []models.Match {
	if len(bestMatches) < 3 || len(bestMatches) == 0 {
		bestMatches = append(bestMatches, match)
		return bestMatches
	}

	// Find which input has the worst matching percent and replace it by the current match if it's higher
	worstMatchIndex := 0
	worstMatchingPercent := 1000.0
	for i := range bestMatches {
		if bestMatches[i].AverageMatchingPercent < worstMatchingPercent {
			worstMatchIndex = i
			worstMatchingPercent = bestMatches[i].AverageMatchingPercent
		}
	}

	if match.AverageMatchingPercent > worstMatchingPercent {
		bestMatches[worstMatchIndex] = match
	}

	return bestMatches
}

func ComputeKeywordsScores2SPLIT(carModel string, keywords []string) []models.KeywordWithScore {
	var keywordsWithScores []models.KeywordWithScore

	carModelKeywords := strings.Split(carModel, " ")

	for i := range keywords {
		bestDistance := 100
		closestKeyword := ""

		bestAverage := 0.0

		distanceFromWholModel := levenshtein.Levenshtein(keywords[i], carModel)
		mp := ComputeMatchingPercentage(len(carModel), distanceFromWholModel)

		//log.Println("--------keyword", keywords[i])
		for j := range carModelKeywords {

			// For each keyword compute Distance between the keyword and the car Model
			distance := levenshtein.Levenshtein(keywords[i], carModelKeywords[j])

			matchingPercentage := ComputeMatchingPercentage(len(carModelKeywords[j]), distance)

			percentMatchingBeginning := ComputeBeginning(keywords[i], carModelKeywords[j])

			weight := 1.0

			if len(keywords[i]) > 1 {
				weight = weight + float64(len(keywords[i]))*0.2
			}

			//newAverage := (matchingPercentage + (percentMatchingBeginning * 1.4)) / 2
			newAverage := ((matchingPercentage * weight) + (percentMatchingBeginning * weight) + mp) / 3
			//log.Printf("matchingPercentage %f weight %f percentMatchingBeginning %f matchPercentFromWholModel %f TOTAL %f", matchingPercentage, weight, percentMatchingBeginning, mp, newAverage)

			//log.Printf("------------carModel %s distance: %d matchingPercentage: %f percentMatchingBeginning: %f average: %f", carModelKeywords[j], distance, matchingPercentage, percentMatchingBeginning, newAverage)

			if newAverage > bestAverage {
				bestDistance = distance
				bestAverage = newAverage
				closestKeyword = carModelKeywords[j]
			}

		}

		// Compute how much percentage of the search keyword is present in the car Model
		//matchingPercentage := ComputeMatchingPercentage(len(closestKeyword), closestDistance)

		keywordWithScore := models.KeywordWithScore{
			Keyword:           keywords[i],
			DistanceFromModel: bestDistance,
			//MatchingPercentage: matchingPercentage,
			MatchingPercentage:       bestAverage,
			ClosestWord:              closestKeyword,
			PercentageFromWholeModel: mp,
		}

		keywordsWithScores = append(keywordsWithScores, keywordWithScore)
	}

	return keywordsWithScores
}

// Compute how much percent of the beginiing characters matches between two strings
func ComputeBeginning(str1, str2 string) float64 {
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
	if distanceFromWord > modelLength {
		return 0
	}

	if modelLength > 0 {
		matchPercent := 100 - (float64(distanceFromWord)*100)/float64(modelLength)
		return math.Round(matchPercent*100) / 100
	}

	return 0
}
