package models

type Match struct {
	Brand                  string
	Model                  string
	KeywordWithScores      []KeywordWithScore
	Distance               int
	AverageMatchingPercent float64
}
