package places

import (
	"blenioviva/internal/db"
	"blenioviva/internal/models"
	"context"
)

// PlaceService holds a reference to the DB
type PlaceService struct {
	DB *db.DB
}

// NewPlaceService returns a new instance of PlaceService
func NewPlaceService(db *db.DB) *PlaceService {
	return &PlaceService{DB: db}
}

func (s *PlaceService) GetByID(ctx context.Context, id string) (*models.Place, error) {
	row := s.DB.GetByID().QueryRow(ctx, `SELECT id, name FROM places WHERE id = $1`, id)

	var place models.Place
	if err := row.Scan(&place.ID, &place.Name); err != nil {
		return nil, err
	}
	return &place, nil
}
