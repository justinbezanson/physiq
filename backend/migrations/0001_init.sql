-- Initial schema for physiq.
-- v1 scope: users, sessions, entries, measurements, goals.
-- TODO: expand as data model milestone lands.

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    weight_unit   TEXT NOT NULL DEFAULT 'lb',
    length_unit   TEXT NOT NULL DEFAULT 'in',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
