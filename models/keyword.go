package models

type KeywordWithScore struct {
	Keyword                  string
	DistanceFromModel        int
	ClosestWord              string
	MatchingPercentage       float64
	PercentageFromWholeModel float64
}
