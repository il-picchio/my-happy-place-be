package models

import "github.com/google/uuid"

type Tag struct {
	ID          uuid.UUID    `json:"id" validate:"required,uuid"`
	Name        string       `json:"name" validate:"required,min=1,max=100"`
	Type        string       `json:"type" validate:"required,oneof=category activity feature"`
	DisplayName Translations `json:"display_name" validate:"required,dive,required"`
	Description Translations `json:"description,omitempty" validate:"omitempty,dive,required"`
}
