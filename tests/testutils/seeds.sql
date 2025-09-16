-- Ensure tables exist (idempotent if already migrated)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name TEXT
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    title TEXT,
    body TEXT
);

-- Seed base data
INSERT INTO users (created_at, updated_at, name) VALUES (now(), now(), 'Alice');
INSERT INTO users (created_at, updated_at, name) VALUES (now(), now(), 'Bob');

INSERT INTO posts (created_at, updated_at, title, body) VALUES (now(), now(), 'Hello', 'World');
INSERT INTO posts (created_at, updated_at, title, body) VALUES (now(), now(), 'Second', 'Post body');


