package matching

import (
	"ghebant/lbc-api/internal/constants"
	"ghebant/lbc-api/models"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestFindBestMatch(t *testing.T) {
	tests := []struct {
		id            int
		search        string
		expectedBrand string
		expectedModel string
	}{
		{0, "Gran Turismo Série5", "BMW", "serie 5"},
		{1, "ds 3 crossback", "Citroen", "ds3"},
		{2, "crossback ds 3", "Citroen", "ds3"},
		{3, "test avant", "Audi", "s4 avant"},
		{4, "s4 avant", "Audi", "s4 avant"},
		{5, "avant s4", "Audi", "s4 avant"},
		{6, "cabr", "Audi", "cabriolet"},
		{7, "s4 av", "Audi", "s4 avant"},
		{8, "s4 cab", "Audi", "s4 cabriolet"},
		{9, "serie5", "BMW", "serie 5"},
		{10, "s4", "Audi", "s4"},
	}

	for i := range tests {
		found := FindBestMatch(tests[i].search, &constants.Vehicles)

		assert.Equal(t, tests[i].expectedBrand, found.Brand, "test %d failed", tests[i].id)
		assert.Equal(t, tests[i].expectedModel, found.Model, "test %d failed", tests[i].id)
	}
}

// TODO Refacto
func TestFindBest(t *testing.T) {
	search := "s4"
	carModel := "s4 avant"
	carModel2 := "s4"

	match := FindBest(search, carModel)
	match2 := FindBest(search, carModel2)
	log.Println(match)
	log.Println(match2)
}

func TestComputeBeginning(t *testing.T) {
	tests := []struct {
		id         int
		str1, str2 string
		expected   float64
	}{
		{0, "s4", "s4 avant", 25},
		{1, "s4", "avant", 0},
		{2, "", "avant", 0},
		{3, "avant", "", 0},
		{4, "avant", "avant", 100},
		{5, "av", "avant", 40},
		{6, "avant", "av", 40},
		{7, "", "", 0},
		{8, "avant", "aavant", 16.666666666666668},
	}

	for i := range tests {
		percentage := ComputeBeginning(tests[i].str1, tests[i].str2)

		assert.Equalf(t, tests[i].expected, percentage, "test %d failed", tests[i].id)
	}
}

func TestComputeKeywordsScores(t *testing.T) {
	tests := []struct {
		carModel string
		keywords []string
		expected []models.KeywordWithScore
	}{
		{
			"ds3",
			[]string{"ds", "3"},
			[]models.KeywordWithScore{
				{"ds", -999, "", -999, 0},
				{"3", -999, "", -999, 0},
			},
		},
	}

	for i := range tests {
		for j, exp := range tests[i].expected {
			keywordsGot := ComputeKeywordsScores(tests[i].carModel, tests[i].keywords)

			assert.Equal(t, len(tests[i].expected), len(keywordsGot))
			assert.Equal(t, exp.Keyword, keywordsGot[j].Keyword)
			assert.Truef(t, keywordsGot[j].DistanceFromModel > 0, "\"%s\" keyword distance is inferior to 0: %d", keywordsGot[j].Keyword, keywordsGot[j].DistanceFromModel)
			assert.Truef(t, keywordsGot[j].MatchingPercentage > 0, "\"%s\" keyword matching is inferior to 0: %f", keywordsGot[j].Keyword, keywordsGot[j].MatchingPercentage)
		}
	}
}

func TestComputeMatchingPercentage(t *testing.T) {
	tests := []struct {
		id          int
		modelLength int
		distance    int
		expected    float64
	}{
		{0, 2, 1, 50},
		{1, 3, 2, 33.33},
		{2, 3, 0, 100},
		{3, 3, 8, 0},
		{4, 0, 0, 0},
	}

	for i := range tests {
		assert.Equalf(t, tests[i].expected, ComputeMatchingPercentage(tests[i].modelLength, tests[i].distance), "test %d failed", tests[i].id)
	}
}

func TestComputeGlobalScores(t *testing.T) {
	tests := []struct {
		id               int
		keywords         []models.KeywordWithScore
		expectedDistance int
		expectedAverage  float64
	}{
		{
			0,
			[]models.KeywordWithScore{
				{"ds", 1, "", 67, 0},
				{"3", 2, "", 34, 0},
				{"crossback", 8, "", 0, 0},
			},
			11,
			33.67,
		},
		{
			1,
			[]models.KeywordWithScore{
				{"test", 4, "", 18, 0},
				{"avant", 4, "", 34, 0},
			},
			8,
			26,
		},
		{
			2,
			[]models.KeywordWithScore{},
			0,
			0,
		},
	}

	for i := range tests {
		globalDistance, globalAverage := ComputeGlobalScores(tests[i].keywords)

		assert.Equalf(t, tests[i].expectedDistance, globalDistance, "test %d failed", tests[i].id)
		assert.Equalf(t, tests[i].expectedAverage, globalAverage, "test %d failed", tests[i].id)
	}
}
