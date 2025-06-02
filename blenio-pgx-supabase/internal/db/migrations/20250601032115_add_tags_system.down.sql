-- Drop indexes first
DROP INDEX IF EXISTS idx_place_tags_place_id;
DROP INDEX IF EXISTS idx_place_tags_tag_id;
DROP INDEX IF EXISTS idx_tags_name;

-- Drop tables in reverse order
DROP TABLE IF EXISTS place_tags;
DROP TABLE IF EXISTS tags;
