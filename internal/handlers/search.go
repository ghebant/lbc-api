package handlers

import (
	"ghebant/lbc-api/internal/levenshtein"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"net/http"
	"strings"
	"unicode"
)

// TODO
// Lowercase everything

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		// TODO add error management
		panic(e)
	}
	return output
}

func NormalizeKeywords(keywords []string) []string {
	var cleanKeywords []string

	for _, keyword := range keywords {
		// lower case
		keyword = strings.ToLower(keyword)
		// remove spaces
		keyword = strings.TrimSpace(keyword)
		// remove accents
		keyword = removeAccents(keyword)

		cleanKeywords = append(cleanKeywords, keyword)
	}

	return cleanKeywords
}

// TEXT = MODEL && PATTERN = INPUT
func AnotherOne(text, pattern string) {

	count := 0
	best := 0

	startIndex := 0
	for i := 0; i < len(text); i++ {
		for j := 0; j < len(pattern); j++ {
			if text[i] != pattern[j] {
				j = 0

				if count > 0 {
					i = startIndex + 1
				}

				if count > best {
					best = count
				}
			}

			if text[i] == pattern[j] {
				if count == 0 {
					startIndex = i
				}

				count += 1
			}
		}
	}

	log.Println("max matching characters:", best)
}

// zabcd
// mabzd
// TODO Utiliser algo Longest Common Substring
func LongestCommonSubstrings(needle, word string) []string {
	var subStrings []string
	count := 0
	matchStartAt := 0

	// wordeas    worldeas

	// https://www.youtube.com/results?search_query=Longest+Common+Substring+golang
	for i := 0; i < len(needle); i++ {
		for j := 0; j < len(word); j++ {
			//log.Printf("needle[%d] %s word[%d] %s count %d", i, string(needle[i]), j, string(word[j]), count)

			if word[j] != needle[i] && count > 0 {
				//log.Println("matchStartAt", matchStartAt, "j", j)
				subStrings = append(subStrings, word[matchStartAt:j])
				//log.Println("found sub string:", word[matchStartAt:j])
				count = 0
			}

			if word[j] == needle[i] {
				if count <= 0 {
					matchStartAt = j
				}
				i += 1
				count += 1
			}

		}
	}

	//log.Println("subs", subStrings)

	return nil
}

func pdr(str1, str2 string) int {
	m := len(str1)
	n := len(str2)

	var LCSuff = make([][]int, m+1)
	for i := range LCSuff {
		LCSuff[i] = make([]int, n+1)
	}

	//subStrings := []string{}

	var result = 0

	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {

			if i == 0 || j == 0 {
				LCSuff[i][j] = 0
			} else if str1[i-1] == str2[j-1] {
				LCSuff[i][j] = LCSuff[i-1][j-1] + 1
				result = max(result, LCSuff[i][j])
				//log.Printf("LCSuff[i][j] %d result %d", LCSuff[i][j], result)
			} else {
				LCSuff[i][j] = 0
			}
		}
	}
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func try(str1, str2 string) int {
	var T = make([][]int, len(str1))
	for i := range T {
		T[i] = make([]int, len(str2))
	}

	var max = 0

	for i := 1; i <= len(str1); i++ {
		for j := 1; j <= len(str2); j++ {
			if str1[i-1] == str2[j-1] {
				T[i][j] = T[i-1][j-1] + 1
				if max < T[i][j] {
					max = T[i][j]
				}
			}
		}
	}
	return max
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

	//Vehicles := map[string][]string{
	//	"Audi":    {"Cabriolet", "Q2", "Q3", "Q5", "Q7", "Q8", "R8", "Rs3", "Rs4", "Rs5", "Rs7", "S3", "S4", "S4 Avant", "S4 Cabriolet", "S5", "S7", "S8", "SQ5", "SQ7", "Tt", "Tts", "V8"},
	//	"BMW":     {"M3", "M4", "M5", "M535", "M6", "M635", "Serie 1", "Serie 2", "Serie 3", "Serie 4", "Serie 5", "Serie 6", "Serie 7", "Serie 8"},
	//	"Citroen": {"C1", "C15", "C2", "C25", "C25D", "C25E", "C25TD", "C3", "C3 Aircross", "C3 Picasso", "C4", "C4 Picasso", "C5", "C6", "C8", "Ds3", "Ds4", "Ds5"},
	//}
	//
	//input := "RS4 avant"
	//
	//if len(input) <= 0 {
	//	log.Println("input len is 0")
	//	return
	//}
	//
	//type found struct {
	//	points int
	//	index  int
	//}
	//
	//var f []found
	//
	//inputWords := strings.Split(input, " ")
	//
	//for index, inputWord := range inputWords {
	//	//point := len(inputWord)
	//
	//	if input == inputWord {
	//		// found
	//		tmp := found{
	//			100, index,
	//		}
	//
	//		f = append(f, tmp)
	//	}
	//}

	c.JSON(http.StatusOK, nil)
}
