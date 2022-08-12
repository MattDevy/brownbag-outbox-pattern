
-- +migrate Up
CREATE SCHEMA IF NOT EXISTS brownbag;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER DATABASE brownbag SET timezone TO 'UTC';
-- +migrate Down
DROP EXTENSION IF EXISTS "uuid-ossp";

DROP SCHEMA IF EXISTS brownbag;