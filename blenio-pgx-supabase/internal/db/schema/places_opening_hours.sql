CREATE TABLE place_opening_hours (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
  day_of_week SMALLINT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6), -- 0 = Monday
  start_time TIME NOT NULL,
  end_time TIME NOT NULL,
  CHECK (start_time < end_time)
);

-- Composite index: fast lookup per place & day
CREATE INDEX idx_opening_hours_place_day ON place_opening_hours (place_id, day_of_week, start_time);
