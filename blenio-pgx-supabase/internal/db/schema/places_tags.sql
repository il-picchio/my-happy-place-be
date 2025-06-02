CREATE TABLE place_tags (
  place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
  tag_id   UUID NOT NULL REFERENCES tags(id)   ON DELETE CASCADE,
  PRIMARY KEY (place_id, tag_id)
);

-- Helpful index for filtering by tag_id
CREATE INDEX idx_place_tags_tag_id ON place_tags(tag_id);
