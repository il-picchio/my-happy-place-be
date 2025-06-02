-- Create the tags table
CREATE TABLE tags (
  id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  name        TEXT        NOT NULL UNIQUE,
  display_name   JSONB    NOT NULL, -- e.g. { "en": "Restaurant", "it": "Ristorante", "de": "Restaurant", "fr": "Restaurant" }
  type        TEXT        NOT NULL
                   CHECK (type IN ('Category','Attribute','Feature')),
  description JSONB       NOT NULL,
  deprecated  BOOLEAN     NOT NULL DEFAULT FALSE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_tags_display_name_gin ON tags USING GIN (display_name jsonb_path_ops);
CREATE INDEX idx_tags_description_gin  ON tags USING GIN (description jsonb_path_ops);

-- Create the place_tags join table
CREATE TABLE place_tags (
  place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
  tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (place_id, tag_id)
);

-- Helpful indexes
CREATE INDEX idx_tags_name ON tags (name);
CREATE INDEX idx_place_tags_tag_id ON place_tags (tag_id);
CREATE INDEX idx_place_tags_place_id ON place_tags (place_id);
