-- name: CreatePlaceTag :exec
INSERT INTO place_tags (
  place_id,
  tag_id
)
VALUES (
  @place_id,  -- place_id UUID
  @tag_id   -- tag_id   UUID
)
ON CONFLICT (place_id, tag_id) DO NOTHING;

-- name: DeletePlaceTag :exec
-- Remove a specific tag assignment from a place
DELETE FROM place_tags
WHERE place_id = @place_id
  AND tag_id   = @tag_id;

  -- db/queries/place_tags.sql

-- name: DeletePlaceTagsExcept :exec
DELETE FROM place_tags
WHERE place_id = @place_id
  AND tag_id <> ALL(@tag_ids);  -- tag_ids UUID[]
