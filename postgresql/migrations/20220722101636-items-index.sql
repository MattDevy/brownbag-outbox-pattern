
-- +migrate Up
CREATE UNIQUE INDEX idx_unique_item_name ON brownbag.items(name);
-- +migrate Down
DROP INDEX brownbag.idx_unique_item_name;
