package levenshtein

// Levenshtein computes the distance from word1 to word2
func Levenshtein(word1, word2 string) int {
	cost := 0
	word1Len := len(word1)
	word2Len := len(word2)
	column := make([]int, len(word1)+1)

	for y := 1; y <= word1Len; y++ {
		column[y] = y
	}

	for x := 1; x <= word2Len; x++ {
		column[0] = x
		lastDiag := x - 1

		for y := 1; y <= word1Len; y++ {
			oldDiag := column[y]
			cost = 0
			if word1[y-1] != word2[x-1] {
				cost = 1
			}

			column[y] = minimum(column[y]+1, column[y-1]+1, lastDiag+cost)
			lastDiag = oldDiag
		}
	}
	return column[word1Len]
}

func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}
