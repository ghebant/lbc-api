package models

type RealEstate struct {
	name string `json:"name"`
}

func (r RealEstate) Name() string {
	return r.name
}

func (r RealEstate) Brand() string {
	return ""
}

func (r RealEstate) Model() string {
	return ""
}
