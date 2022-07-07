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
		{0, "rs4 avant", "Audi", "rs4"},
		{1, "Gran Turismo Série5", "BMW", "serie 5"},
		{2, "ds 3 crossback", "Citroen", "ds3"},
		{3, "crossback ds 3", "Citroen", "ds3"},
		{4, "test avant", "Audi", "s4 avant"},
		{5, "cabr", "Audi", "cabriolet"},
		{6, "s4 avant", "Audi", "s4 avant"},
		// TODO Donner point si sub present au debut
		//{"s4 c", "Audi", "s4 cabriolet"},
	}

	for i := range tests {
		found := FindBestMatch(tests[i].search, &constants.Vehicles)

		assert.Equal(t, tests[i].expectedBrand, found.Brand, "test %d failed", tests[i].id)
		assert.Equal(t, tests[i].expectedModel, found.Model, "test %d failed", tests[i].id)
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
				{"ds", -999, -999},
				{"3", -999, -999},
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
		{3, 3, 8, -166.67},
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
				{"ds", 1, 67},
				{"3", 2, 34},
				{"crossback", 8, 0},
			},
			11,
			33.67,
		},
		{
			1,
			[]models.KeywordWithScore{
				{"test", 4, 18},
				{"avant", 4, 34},
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
