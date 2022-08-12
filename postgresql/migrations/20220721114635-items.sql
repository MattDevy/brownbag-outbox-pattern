
-- +migrate Up
CREATE TABLE IF NOT EXISTS brownbag.items (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name varchar NOT NULL,
    count int default 0,
    price float(2)
);
-- +migrate Down
DROP TABLE IF EXISTS brownbag.items;