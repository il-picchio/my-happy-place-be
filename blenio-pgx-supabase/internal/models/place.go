package models

import "github.com/google/uuid"

type Place struct {
	ID           uuid.UUID    `json:"id" validate:"required,uuid"`
	Title        Translations `json:"title"`
	Description  Translations `json:"description"`
	Geo          Geo          `json:"geo" validate:"required"`
	Address      Address      `json:"address" validate:"required"`
	PhotoURLs    []string     `json:"photo_urls,omitempty" validate:"dive,url"`
	OpeningHours OpeningHours `json:"opening_hours,omitempty" validate:"omitempty"`
	Tags         []Tag        `json:"tags,omitempty" validate:"omitempty,dive,required"`
}

type GetPlacesNearbyAfterRequest struct {
	Lat          float64  `json:"lat" validate:"required,gte=-90,lte=90"`
	Lng          float64  `json:"lng" validate:"required,gte=-180,lte=180"`
	LastDistance *float64 `json:"last_distance,omitempty" validate:"omitempty,gte=0"`
	MaxDistance  *float64 `json:"max_distance,omitempty" validate:"omitempty,gte=0"`
}

type PlaceWithDistance struct {
	Place            // embedded → all Place fields included
	Distance float64 `json:"distance"`
}

type GetPlaceOptions struct {
	WithOpeningHours bool
	// add more options later
}

type UpdatePlacePartialParams struct {
	ID          uuid.UUID     `validate:"required"`
	Title       *Translations `validate:"omitempty,dive,required"`
	Description *Translations `validate:"omitempty,dive,required"`
	Lat         *float64      `validate:"omitempty,gte=-90,lte=90"`
	Lng         *float64      `validate:"omitempty,gte=-180,lte=180"`
	Street      *string       `validate:"omitempty,min=1"`
	Zip         *string       `validate:"omitempty,ch_zip"`
	City        *string       `validate:"omitempty,min=1"`
	State       *string       `validate:"omitempty,min=1"`
	Country     *string       `validate:"omitempty,country"`
	PhotoUrls   *[]string     `validate:"omitempty,dive,url"`
}
