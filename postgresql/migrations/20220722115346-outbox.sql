
-- +migrate Up

CREATE SEQUENCE IF NOT EXISTS outbox_events_offset_seq AS INTEGER START 1 INCREMENT BY 1;
CREATE TABLE IF NOT EXISTS brownbag.outbox_events (
    "offset" INTEGER NOT NULL DEFAULT nextval('outbox_events_offset_seq'),
    "uuid" VARCHAR(36) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "payload" JSON DEFAULT NULL,
    "metadata" JSON DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS brownbag.outbox_offsets_events (
    consumer_group VARCHAR(255) NOT NULL,
    offset_acked BIGINT,
    offset_consumed BIGINT NOT NULL,
    PRIMARY KEY(consumer_group)
);
-- +migrate Down
DROP TABLE IF EXISTS brownbag.outbox_events;
DROP TABLE IF EXISTS brownbag.outbox_offset_events;
