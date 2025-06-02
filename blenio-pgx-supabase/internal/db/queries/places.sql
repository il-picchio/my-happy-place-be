-- name: CreatePlace :exec
INSERT INTO places (
  id,
  geo,
  title,
  description,
  street,
  zip,
  city,
  state,
  country,
  photo_urls,
  created_at,
  updated_at
)
VALUES (
  @id,
  ST_MakePoint(@lng::float8, @lat::float8),
  @title,
  @description,
  @street,
  @zip,
  @city,
  @state,
  @country,
  @photo_urls,
  now(), now()
);

-- name: GetPlaceByID :one
SELECT
  id,
  title,
  description,
  ST_X(geo::geometry)::float8 AS lng,
  ST_Y(geo::geometry)::float8 AS lat,
  street,
  zip,
  city,
  state,
  country,
  photo_urls
FROM places
WHERE id = @id;

-- name: GetPlacesNearbyWithCoords :many
-- Find places near a given point, using distance sorting + cursor pagination,
-- filtered by matching all provided tags and by opening hours.
-- Parameters:
--   @lng            : float8    — reference longitude
--   @lat            : float8    — reference latitude
--   @last_distance  : float8    — distance cursor for pagination (pass NULL for first page)
--   @max_distance   : float8    — maximum distance filter (meters) (pass NULL to ignore)
--   @tag_ids        : uuid[]    — array of tag IDs that the place must have (pass NULL or '{}' to ignore)
--   @day_of_week    : int4      — integer 0=Sunday … 6=Saturday (pass NULL to ignore)
--   @time_of_day    : time      — time in HH:MM:SS (pass NULL to ignore)
WITH nearest AS (
  SELECT
    p.id,
    p.title,
    p.description,
    ST_X(p.geo::geometry)::float8 AS lng,
    ST_Y(p.geo::geometry)::float8 AS lat,
    p.street,
    p.zip,
    p.city,
    p.state,
    p.country,
    p.photo_urls,
    p.geo
  FROM places AS p
  ORDER BY p.geo::geometry <-> ST_SetSRID(ST_MakePoint(@lng::float8, @lat::float8), 4326)
  LIMIT 200
)
SELECT
  n.id,
  n.title,
  n.description,
  n.lng,
  n.lat,
  n.street,
  n.zip,
  n.city,
  n.state,
  n.country,
  n.photo_urls,
  ST_Distance(n.geo, ST_SetSRID(ST_MakePoint(@lng::float8, @lat::float8), 4326))::float8 AS distance
FROM nearest AS n
WHERE
  -- Cursor pagination by distance
  (
    @last_distance IS NULL
    OR ST_Distance(n.geo, ST_SetSRID(ST_MakePoint(@lng::float8, @lat::float8), 4326)) > @last_distance
  )
  AND
  -- Max distance filter
  (
    @max_distance IS NULL
    OR ST_Distance(n.geo, ST_SetSRID(ST_MakePoint(@lng::float8, @lat::float8), 4326)) <= @max_distance
  )
  AND
  -- Tag filtering: require every tag_id in @tag_ids
  (
    @tag_ids IS NULL
    OR @tag_ids = '{}'
    OR NOT EXISTS (
      SELECT required.tag_id
      FROM unnest(@tag_ids) AS required(tag_id)
      WHERE NOT EXISTS (
        SELECT 1
        FROM place_tags AS pt
        WHERE pt.place_id = n.id
          AND pt.tag_id = required.tag_id
      )
    )
  )
  AND
  -- Opening hours filter
  (
    @day_of_week IS NULL
    OR @time_of_day IS NULL
    OR EXISTS (
      SELECT 1
      FROM place_opening_hours AS oh
      WHERE oh.place_id = n.id
        AND oh.day_of_week = @day_of_week
        AND @time_of_day >= oh.open_time
        AND @time_of_day <  oh.close_time
    )
  )
ORDER BY distance ASC, id ASC
LIMIT 20;


-- name: GetPlacesFilteredNoCoords :many
-- Fetch places ordered by id + keyset pagination, filtered by tags and opening hours.
-- Parameters:
--   @last_id        : uuid     — last place.id from previous page (pass NULL for first page)
--   @tag_ids        : uuid[]   — array of tag IDs that the place must have (pass NULL or '{}' to ignore)
--   @day_of_week    : int4     — integer 0=Sunday … 6=Saturday (pass NULL to ignore)
--   @time_of_day    : time     — time in HH:MM:SS (pass NULL to ignore)
SELECT
  p.id,
  p.title,
  p.description,
  ST_X(p.geo::geometry)::float8 AS lng,
  ST_Y(p.geo::geometry)::float8 AS lat,
  p.street,
  p.zip,
  p.city,
  p.state,
  p.country,
  p.photo_urls,
  NULL::float8 AS distance
FROM places AS p
WHERE
  -- Keyset pagination by id
  (
    @last_id IS NULL
    OR p.id > @last_id
  )
  AND
  -- Tag filtering: require every tag_id in @tag_ids
  (
    @tag_ids IS NULL
    OR @tag_ids = '{}'
    OR NOT EXISTS (
      SELECT required.tag_id
      FROM unnest(@tag_ids) AS required(tag_id)
      WHERE NOT EXISTS (
        SELECT 1
        FROM place_tags AS pt
        WHERE pt.place_id = p.id
          AND pt.tag_id = required.tag_id
      )
    )
  )
  AND
  -- Opening hours filter
  (
    @day_of_week IS NULL
    OR @time_of_day IS NULL
    OR EXISTS (
      SELECT 1
      FROM place_opening_hours AS oh
      WHERE oh.place_id = p.id
        AND oh.day_of_week = @day_of_week
        AND @time_of_day >= oh.open_time
        AND @time_of_day <  oh.close_time
    )
  )
ORDER BY p.id ASC
LIMIT 20;


-- name: UpdatePlacePartial :exec
UPDATE places
SET
  title = COALESCE(@title, title),
  description = COALESCE(@description, description),
  geo = CASE
          WHEN @lng IS NOT NULL AND @lat IS NOT NULL
          THEN ST_SetSRID(ST_MakePoint(@lng::float8, @lat::float8), 4326)
          ELSE geo
        END,
  street = COALESCE(@street, street),
  zip = COALESCE(@zip, zip),
  city = COALESCE(@city, city),
  state = COALESCE(@state, state),
  country = COALESCE(@country, country),
  photo_urls = COALESCE(@photo_urls, photo_urls),
  updated_at = now()
WHERE id = @id;





