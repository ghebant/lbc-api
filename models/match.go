package models

type Match struct {
	Brand                  string
	Model                  string
	KeywordWithScores      []KeywordWithScore
	AverageMatchingPercent float64
}
