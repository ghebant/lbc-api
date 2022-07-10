package matching

import (
	"ghebant/lbc-api/internal/constants"
	"ghebant/lbc-api/models"
	"github.com/stretchr/testify/assert"
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

func TestComputeMatch(t *testing.T) {
	tests := []struct {
		id              int
		search          string
		model           string
		expectedAverage float64
	}{
		{0, "Gran Turismo Série5", "serie 5", 18.86},
		{1, "ds 3 crossback", "ds3", 35.56},
		{2, "crossback ds 3", "ds3", 35.56},
		{3, "test avant", "s4 avant", 85.17},
		{4, "s4 avant", "s4 avant", 1000},
		{5, "avant s4", "s4 avant", 127.92},
		{6, "cabr", "cabriolet", 68.14},
		{7, "s4 av", "s4 avant", 73.67},
		{8, "s4 cab", "s4 cabriolet", 71.39},
		{9, "serie5", "serie 5", 148.35},
		{10, "s4", "s4", 1000},
	}

	for i := range tests {
		match := ComputeMatch(tests[i].search, tests[i].model)

		assert.Equal(t, tests[i].expectedAverage, match.AverageMatchingPercent, "test %d failed", tests[i].id)
	}
}

func TestComputeCharactersMatching(t *testing.T) {
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
		percentage := ComputeCharactersMatching(tests[i].str1, tests[i].str2)

		assert.Equalf(t, tests[i].expected, percentage, "test %d failed", tests[i].id)
	}
}

func TestComputeKeywordsScores(t *testing.T) {
	tests := []struct {
		id       int
		carModel string
		keywords []string
		expected []models.KeywordWithScore
	}{
		{
			0,
			"ds3",
			[]string{"ds", "3"},
			[]models.KeywordWithScore{
				{"ds", 84.44711111111111},
				{"3", 22.22},
			},
		},
	}

	for i := range tests {
		for j, exp := range tests[i].expected {
			keywordsGot := ComputeKeywordsScores(tests[i].carModel, tests[i].keywords)

			assert.Equal(t, len(tests[i].expected), len(keywordsGot))
			assert.Equal(t, exp.Keyword, keywordsGot[j].Keyword)
			assert.Equalf(t, exp.MatchingPercentage, keywordsGot[j].MatchingPercentage, "test %d failed", tests[i].id)
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
				{"ds", 67},
				{"3", 34},
				{"crossback", 0},
			},
			11,
			33.67,
		},
		{
			1,
			[]models.KeywordWithScore{
				{"test", 18},
				{"avant", 34},
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
		globalAverage := ComputeGlobalScores(tests[i].keywords)

		assert.Equalf(t, tests[i].expectedAverage, globalAverage, "test %d failed", tests[i].id)
	}
}
