-- Query to insert a new tag and return all columns
-- name: CreateTag :one
INSERT INTO tags (name, display_name, type, description)
VALUES (@name, @displayname, @type, @description)
RETURNING id, name, display_name, type, description, deprecated, created_at, updated_at;

-- name: GetAllTagsByType :many
-- Fetch all tags of a given type, optionally including deprecated ones.
SELECT
  id,
  name,
  display_name,
  type,
  description,
  deprecated,
  created_at,
  updated_at
FROM tags
WHERE type = @type
  AND (deprecated = FALSE OR @includeDeprecated::boolean)
ORDER BY name;

-- name: GetTagsLike :many
-- Search tags by partial name, any display_name value, or any description value (across all languages),
-- optionally including deprecated ones.
SELECT
  id,
  name,
  display_name,
  type,
  description,
  deprecated,
  created_at,
  updated_at
FROM tags
WHERE (
      name ILIKE '%' || @term || '%'
   OR EXISTS (
         SELECT 1
         FROM jsonb_each_text(display_name) AS kv(key, value)
         WHERE value ILIKE '%' || @term || '%'
      )
   OR EXISTS (
         SELECT 1
         FROM jsonb_each_text(description) AS kv(key, value)
         WHERE value ILIKE '%' || @term || '%'
      )
  )
  AND (deprecated = FALSE OR @includeDeprecated::boolean)
ORDER BY
  (name ILIKE @term || '%') DESC,
  name;

-- name: GetTagByExactName :one
-- Fetch a single tag by exact name, optionally including deprecated.
SELECT
  id,
  name,
  display_name,
  type,
  description,
  deprecated,
  created_at,
  updated_at
FROM tags
WHERE name = @name
  AND (deprecated = FALSE OR @includeDeprecated::boolean);

-- name: DeleteTag :exec
-- Permanently delete a tag by its ID.
DELETE FROM tags
WHERE id = @id;

-- name: DeprecateTag :one
-- Mark a tag as deprecated.
UPDATE tags
SET deprecated = TRUE,
    updated_at = now()
WHERE id = @id
RETURNING
  id,
  name,
  display_name,
  type,
  description,
  deprecated,
  created_at,
  updated_at;

  SELECT
  t.id,
  t.name,
  t.display_name,
  t.type,
  t.description,
  t.deprecated,
  t.created_at,
  t.updated_at
FROM tags AS t
JOIN place_tags AS pt
  ON t.id = pt.tag_id
WHERE pt.place_id = @place_id
ORDER BY t.name;

-- name: UpdateTag :one
-- Update only the provided fields for a specific tag; leave others unchanged.
UPDATE tags
SET
  name = COALESCE(@name, name),
  display_name = COALESCE(@display_name, display_name),
  type = COALESCE(@type, type),
  description = COALESCE(@description, description),
  deprecated = COALESCE(@deprecated, deprecated),
  updated_at = now()
WHERE id = @id
RETURNING
  id,
  name,
  display_name,
  type,
  description,
  deprecated,
  created_at,
  updated_at;
