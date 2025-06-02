package db

import (
	"blenioviva/internal/db/generated"
	pgxutils "blenioviva/internal/db/utils"
	"blenioviva/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (db *DB) GetOpeningHours(ctx context.Context, placeID uuid.UUID) (*models.OpeningHours, error) {
	rows, err := db.queries.GetOpeningHoursByPlaceID(ctx, placeID)
	if err != nil {
		return nil, err
	}

	weekly := make(map[int8]models.DaySchedule)

	for _, row := range rows {
		day := int8(row.DayOfWeek)

		start, err1 := row.StartTime.Value()
		end, err2 := row.EndTime.Value()

		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid time value for place %s on day %d: start=%v, end=%v", placeID, day, row.StartTime, row.EndTime)
		}

		schedule := weekly[day]
		schedule.Intervals = append(schedule.Intervals, models.TimeRange{
			Start: start.(time.Time).Format("15:04"),
			End:   end.(time.Time).Format("15:04"),
		})
		weekly[day] = schedule
	}

	return &models.OpeningHours{
		Weekly: weekly,
	}, nil
}

func (db *DB) InsertWeeklyOpeningHours(
	ctx context.Context,
	tx pgx.Tx, // not *pgx.Tx — allows both pool and tx interfaces
	placeID uuid.UUID,
	opening map[int8]models.DaySchedule,
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

	for day, schedule := range opening {
		for _, interval := range schedule.Intervals {
			start, err := pgxutils.ParsePgTime(interval.Start)
			if err != nil {
				return fmt.Errorf("invalid start time %q: %w", interval.Start, err)
			}
			end, err := pgxutils.ParsePgTime(interval.End)
			if err != nil {
				return fmt.Errorf("invalid end time %q: %w", interval.End, err)
			}

			err = qtx.CreateOpeningHour(ctx, generated.CreateOpeningHourParams{
				PlaceID:   placeID,
				DayOfWeek: int16(day),
				StartTime: start,
				EndTime:   end,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *DB) UpdateOpeningHoursForDay(
	ctx context.Context,
	tx pgx.Tx, // optional
	placeID uuid.UUID,
	day int8,
	schedule models.DaySchedule,
) error {
	var err error

	newTx := false
	if tx == nil {
		newTx = true
		tx, err = db.pool.Begin(ctx)
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				_ = tx.Rollback(ctx)
			}
		}()
	}

	qtx := db.queries.WithTx(tx)

	// 1. Delete existing for that day
	if err = qtx.DeleteOpeningHoursByDay(ctx, generated.DeleteOpeningHoursByDayParams{
		PlaceID:   placeID,
		DayOfWeek: int16(day),
	}); err != nil {
		return err
	}

	// 2. Insert all new intervals
	for _, interval := range schedule.Intervals {
		start, err := pgxutils.ParsePgTime(interval.Start)
		if err != nil {
			return fmt.Errorf("invalid start time %q: %w", interval.Start, err)
		}
		end, err := pgxutils.ParsePgTime(interval.End)
		if err != nil {
			return fmt.Errorf("invalid end time %q: %w", interval.End, err)
		}

		err = qtx.CreateOpeningHour(ctx, generated.CreateOpeningHourParams{
			PlaceID:   placeID,
			DayOfWeek: int16(day),
			StartTime: start,
			EndTime:   end,
		})
		if err != nil {
			return err
		}
	}

	if newTx {
		err = tx.Commit(ctx)
	}
	return err
}
