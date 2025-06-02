CREATE TABLE places (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  geo GEOGRAPHY(Point, 4326) NOT NULL,
  title JSONB NOT NULL,
  description JSONB,
  street TEXT NOT NULL,
  zip TEXT NOT NULL,
  city TEXT NOT NULL,
  state TEXT NOT NULL,
  country TEXT NOT NULL,
  photo_urls TEXT[],
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_places_geo ON places USING GIST (geo);
