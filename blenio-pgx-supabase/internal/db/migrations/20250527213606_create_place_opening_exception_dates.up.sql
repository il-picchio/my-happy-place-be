CREATE TABLE place_opening_exception_dates (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	place_id UUID REFERENCES places(id) ON DELETE CASCADE,
	date DATE NOT NULL,
	closed BOOLEAN DEFAULT FALSE,
	reason JSONB DEFAULT '{}'::JSONB, -- e.g. { "en": "Public holiday", "rm": "Festa publica" }
);
