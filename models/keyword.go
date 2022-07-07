package models

type KeywordWithScore struct {
	Keyword            string
	DistanceFromModel  int
	MatchingPercentage float64
}
