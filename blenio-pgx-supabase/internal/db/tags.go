package db

import (
	"blenioviva/internal/db/generated"
	"blenioviva/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (db *DB) InsertPlaceTags(
	ctx context.Context,
	tx pgx.Tx, // allows using either a transaction or direct pool
	placeID uuid.UUID,
	tagIDs []uuid.UUID,
) (err error) {
	ownTx := false

	// Start a new transaction if none was provided
	if tx == nil {
		var newTx pgx.Tx
		newTx, err = db.pool.Begin(ctx)
		if err != nil {
			return err
		}
		tx = newTx
		ownTx = true
	}

	if ownTx {
		defer func() {
			if err != nil {
				_ = tx.Rollback(ctx)
			} else {
				err = tx.Commit(ctx)
			}
		}()
	}

	qtx := db.queries.WithTx(tx)

	for _, tagID := range tagIDs {
		err = qtx.CreatePlaceTag(ctx, generated.CreatePlaceTagParams{
			PlaceID: placeID,
			TagID:   tagID,
		})
		if err != nil {
			return fmt.Errorf("failed to assign tag %s to place %s: %w", tagID, placeID, err)
		}
	}

	return nil
}

// CreateTag creates a new tag in the database. If tx is nil, it starts its own transaction.
// Parameters:
//   - ctx:       the context
//   - tx:        an existing pgx.Tx (or nil to begin a new one)
//   - name:      unique tag name (e.g. "veg-friendly")
//   - tagType:   one of "Category", "Attribute", or "Feature"
//   - displayName: map[string]string with localized labels, e.g. {"en":"Vegan Friendly", "it":"Vegano"}
//   - description: map[string]string with localized descriptions
//
// Returns the created Tag record or an error.
func (db *DB) CreateTag(
	ctx context.Context,
	tx pgx.Tx, // nil to begin a new TX
	name string,
	tagType string,
	displayName map[string]string,
	description map[string]string,
) (models.Tag, error) {
	var result models.Tag
	ownTx := false

	// 1) Start a TX if none provided
	if tx == nil {
		newTx, err := db.pool.Begin(ctx)
		if err != nil {
			return result, err
		}
		tx = newTx
		ownTx = true
	}

	// 2) Commit or rollback if we started the TX
	if ownTx {
		defer func() {
			if result.ID == uuid.Nil { // indicates error
				_ = tx.Rollback(ctx)
			} else {
				_ = tx.Commit(ctx)
			}
		}()
	}

	qtx := db.queries.WithTx(tx)

	// 3) Call the SQLC‐generated CreateTag
	row, err := qtx.CreateTag(ctx, generated.CreateTagParams{
		Name:        name,
		Type:        tagType,
		Displayname: displayName,
		Description: description,
	})
	if err != nil {
		return result, fmt.Errorf("CreateTag(%s): %w", name, err)
	}

	// 4) Map generated.CreateTagRow → models.Tag
	result = models.Tag{
		ID:          row.ID,
		Name:        row.Name,
		Type:        row.Type,
		DisplayName: row.DisplayName,
		Description: row.Description,
	}

	return result, nil
}
