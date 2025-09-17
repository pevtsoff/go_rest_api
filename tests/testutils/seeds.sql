-- Seed base data
INSERT INTO users (created_at, updated_at, name) VALUES (now(), now(), 'Alice');
INSERT INTO users (created_at, updated_at, name) VALUES (now(), now(), 'Bob');

INSERT INTO posts (created_at, updated_at, title, body) VALUES (now(), now(), 'Hello', 'World');
INSERT INTO posts (created_at, updated_at, title, body) VALUES (now(), now(), 'Second', 'Post body');


