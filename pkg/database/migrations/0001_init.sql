-- +goose Up
CREATE TABLE IF NOT EXISTS parts (
	object_id TEXT NOT NULL, 
	node_id INT DEFAULT 0, 
	part_id INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS objects (
	object_id TEXT NOT NULL, 
	size INT DEFAULT 0
);

CREATE INDEX IF NOT EXISTS parts_idx ON parts (object_id);
CREATE INDEX IF NOT EXISTS objects_idx ON parts (object_id);

-- +goose Down
DROP TABLE parts;
DROP INDEX parts_idx;
DROP INDEX objects_idx;
