package db

import (
	"blenioviva/internal/db/generated"
	"blenioviva/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (db *DB) GetPlaceByID(ctx context.Context, id string, opts *models.GetPlaceOptions) (*models.Place, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	row, err := db.queries.GetPlaceByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	place := &models.Place{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		Geo: models.Geo{
			Lat: row.Lat,
			Lng: row.Lng,
		},
		Address: models.Address{
			Street:  row.Street,
			Zip:     row.Zip,
			City:    row.City,
			State:   row.State,
			Country: row.Country,
		},
		PhotoURLs: row.PhotoUrls,
	}

	// Fetch opening hours if requested
	if opts != nil && opts.WithOpeningHours {
		opening, err := db.GetOpeningHours(ctx, parsedID)
		if err != nil {
			return nil, err
		}
		place.OpeningHours = *opening
	}

	return place, nil
}

func (db *DB) GetPlacesNearbyAfter(
	ctx context.Context,
	lat, lng float64,
	lastDistance *float64, // pass nil for first page
	maxDistance *float64, // pass nil to ignore
	tagIds *[]uuid.UUID,
	opts *models.GetPlaceOptions,
) ([]models.PlaceWithDistance, error) {
	params := generated.GetPlacesNearbyWithCoordsParams{
		Lat:          lat,
		Lng:          lng,
		LastDistance: lastDistance,
		MaxDistance:  maxDistance,
		TagIds:       tagIds,
	}

	rows, err := db.queries.GetPlacesNearbyWithCoords(ctx, params)
	if err != nil {
		return nil, err
	}

	places := make([]models.PlaceWithDistance, len(rows))
	for i, row := range rows {
		places[i] = models.PlaceWithDistance{
			Place: models.Place{
				ID:          row.ID,
				Title:       row.Title,
				Description: row.Description,
				Geo: models.Geo{
					Lat: row.Lat,
					Lng: row.Lng,
				},
				Address: models.Address{
					Street:  row.Street,
					Zip:     row.Zip,
					City:    row.City,
					State:   row.State,
					Country: row.Country,
				},
				PhotoURLs: row.PhotoUrls,
			},
			Distance: row.Distance,
		}
	}

	return places, nil
}

func (db *DB) CreatePlace(ctx context.Context, place *models.Place) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err = db.insertPlace(ctx, tx, place); err != nil {
		return err
	}

	if place.OpeningHours.Weekly != nil {
		if err = db.InsertWeeklyOpeningHours(ctx, tx, place.ID, place.OpeningHours.Weekly); err != nil {
			return err
		}
	}

	tagIDs := make([]uuid.UUID, len(place.Tags))
	for i, t := range place.Tags {
		tagIDs[i] = t.ID
	}

	if len(tagIDs) > 0 {
		if err = db.InsertPlaceTags(ctx, tx, place.ID, tagIDs); err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	return err
}

func (db *DB) UpdatePlacePartial(
	ctx context.Context,
	params generated.UpdatePlacePartialParams,
) error {
	// ✅ Apply manual logic for interdependent fields (e.g. Lat/Lng)
	if (params.Lat != nil && params.Lng == nil) || (params.Lat == nil && params.Lng != nil) {
		return fmt.Errorf("both lat and lng must be provided together")
	}

	// ✅ Perform the update
	return db.queries.UpdatePlacePartial(ctx, params)
}

func (db *DB) insertPlace(ctx context.Context, tx pgx.Tx, place *models.Place) error {
	qtx := db.queries.WithTx(tx)
	params := generated.CreatePlaceParams{
		ID:          place.ID,
		Lng:         place.Geo.Lng,
		Lat:         place.Geo.Lat,
		Title:       place.Title,
		Description: place.Description,
		Street:      place.Address.Street,
		Zip:         place.Address.Zip,
		City:        place.Address.City,
		State:       place.Address.State,
		Country:     place.Address.Country,
		PhotoUrls:   place.PhotoURLs,
	}
	return qtx.CreatePlace(ctx, params)
}
