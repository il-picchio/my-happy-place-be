-- name: CreateOpeningHour :exec
INSERT INTO place_opening_hours (
  place_id,
  day_of_week,
  start_time,
  end_time
) VALUES (
  @place_id,
  @day_of_week,
  @start_time,
  @end_time
);

-- name: GetOpeningHoursByPlaceID :many
SELECT
  day_of_week,
  start_time,
  end_time
FROM place_opening_hours
WHERE place_id = @place_id
ORDER BY day_of_week, start_time;

-- name: DeleteOpeningHoursByDay :exec
DELETE FROM place_opening_hours
WHERE place_id = @place_id AND day_of_week = @day_of_week;
