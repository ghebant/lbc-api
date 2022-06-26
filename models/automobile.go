package models

type Automobile struct {
	name  string `json:"name"`
	brand string `json:"brand"`
	model string `json:"model"`
}

func (a Automobile) Name() string {
	return a.name
}

func (a Automobile) Brand() string {
	return a.brand
}

func (a Automobile) Model() string {
	return a.model
}
