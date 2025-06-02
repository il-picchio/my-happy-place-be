package models

type User struct {
	ID          string   `json:"id" validate:"required,uuid"`
	Email       string   `json:"email" validate:"required,email"`
	Name        string   `json:"name" validate:"required"`
	Companies   []string `json:"companies" validate:"dive,required"`
	PhoneNumber *string  `json:"phone_number" validate:"omitempty,e164"`
	PhotoURL    *string  `json:"photo_url" validate:"omitempty,url"`
}
