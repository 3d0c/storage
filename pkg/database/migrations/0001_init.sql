-- +goose Up
CREATE TABLE IF NOT EXISTS parts (
	object_id TEXT NOT NULL, 
	node_id INT DEFAULT 0, 
	part_id INT DEFAULT 0
);

-- @TODO. Check whether is it possible to add integer range to the index
CREATE INDEX IF NOT EXISTS parts_idx ON parts (object_id);

-- +goose Down
DROP TABLE parts;
DROP INDEX parts_idx;
