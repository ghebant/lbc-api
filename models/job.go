package models

type Job struct {
	name string `json:"name"`
}

func (j Job) Name() string {
	return j.name
}

func (j Job) Brand() string {
	return ""
}

func (j Job) Model() string {
	return ""
}
