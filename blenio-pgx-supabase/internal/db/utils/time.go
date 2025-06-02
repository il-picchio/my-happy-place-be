package pgxutils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ParsePgTime converts a "15:04" string into pgtype.Time
func ParsePgTime(s string) (pgtype.Time, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return pgtype.Time{}, err
	}

	var pgTime pgtype.Time
	if err := pgTime.Scan(t); err != nil {
		return pgtype.Time{}, err
	}

	return pgTime, nil
}
