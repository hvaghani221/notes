-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_notes_content ON notes USING GIN (to_tsvector('english', content));

-- +goose Down
DROP INDEX IF EXISTS idx_notes_content;
DROP EXTENSION IF EXISTS pg_trgm;

