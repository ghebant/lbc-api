package models

type Category interface {
	Name() string

	// Automobile
	Brand() string
	Model() string
}
