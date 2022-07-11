package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveAccents(t *testing.T) {
	tests := []struct {
		id       int
		input    string
		expected string
	}{
		{0, "Bien joué à", "Bien joue a"},
		{1, "", ""},
	}

	for i := range tests {
		normalizedStr := RemoveAccents(tests[i].input)

		assert.Equalf(t, tests[i].expected, normalizedStr, "test %d failed", tests[i].id)
	}
}

func TestNormalizeKeywords(t *testing.T) {
	tests := []struct {
		id       int
		input    string
		expected string
	}{
		{0, "Bien joué à", "bien joue a"},
		{1, "", ""},
	}

	for i := range tests {
		normalizedStr := NormalizeString(tests[i].input)

		assert.Equalf(t, tests[i].expected, normalizedStr, "test %d failed", tests[i].id)
	}
}
