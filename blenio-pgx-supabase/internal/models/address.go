package models

type Address struct {
	Place   Translations `json:"place" validate:"required"`
	Street  string       `json:"street" validate:"required"`
	City    string       `json:"city" validate:"required"`
	State   string       `json:"state" validate:"required"`
	Zip     string       `json:"zip" validate:"required,ch_zip"`
	Country string       `json:"country" validate:"required,country"`
}
