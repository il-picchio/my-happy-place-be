CREATE TABLE place_opening_exception_hours (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	exception_date_id UUID REFERENCES place_opening_exception_dates(id) ON DELETE CASCADE,
	start_time TIME NOT NULL,
	end_time TIME NOT NULL
);

CREATE INDEX idx_exception_hours_exception_id ON place_opening_exception_hours(exception_date_id);
