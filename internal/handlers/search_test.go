package handlers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveAccents(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Bien joué à", "Bien joue a"},
		{"", ""},
	}

	for i := range tests {
		normalizedStr := RemoveAccents(tests[i].input)

		assert.Equal(t, tests[i].expected, normalizedStr)
	}
}

func TestNormalizeKeywords(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Bien joué à", "bien joue a"},
		{"", ""},
	}

	for i := range tests {
		normalizedStr := NormalizeString(tests[i].input)

		assert.Equal(t, tests[i].expected, normalizedStr)
	}
}

func TestSearchCar(t *testing.T) {
	tests := []struct {
		search        string
		expectedBrand string
		expectedModel string
	}{
		{"rs4 avant", "Audi", "rs4"},
		{"Gran Turismo Série5", "BMW", "serie 5"},
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
		assert.Equal(t, tests[i].expected, ComputeMatchingPercentage(tests[i].modelLength, tests[i].distance))
	}
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
		tests[i].bestMatches = Top3BestMatches(tests[i].match, tests[i].bestMatches)

		assert.Equal(t, tests[i].expected, tests[i].bestMatches)
	}
}
