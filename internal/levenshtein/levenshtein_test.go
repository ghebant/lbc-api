package levenshtein

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		id               int
		search           string
		match            string
		expectedDistance int
	}{
		{0, "c4 picasso", "c4 picasso", 0},
		{1, "3", "3", 0},
		{2, "aircross c3", "c3 aircross", 6},
		{3, "ds", "ds3", 1},
		{4, "3", "ds3", 2},
		{5, "avant", "s4 avant", 3},
		{6, "", "s4 avant", 8},
		{7, "abc", "", 3},
		{8, "", "", 0},
	}

	for i := range tests {
		distance := Levenshtein(tests[i].search, tests[i].match)

		assert.Equalf(t, tests[i].expectedDistance, distance, "test %d failed", tests[i].id)
	}
}

func TestMinimum(t *testing.T) {
	tests := []struct {
		id, a, b, c     int
		expectedMinimum int
	}{
		{0, 1, 2, 3, 1},
		{1, 18, 0, 3, 0},
		{2, 1, 1, -1, -1},
		{3, 0, 0, 0, 0},
	}

	for i := range tests {
		min := minimum(tests[i].a, tests[i].b, tests[i].c)

		assert.Equalf(t, tests[i].expectedMinimum, min, "test %d failed", tests[i].id)
	}
}
