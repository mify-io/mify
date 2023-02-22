-- Check out sqlc docs to get started:
-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html
-- Here are simple CRUD examples:

-- name: GetAuthor :one
-- SELECT * FROM authors
-- WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
-- SELECT * FROM authors
-- ORDER BY name;

-- name: CreateAuthor :one
-- INSERT INTO authors (
--   name, bio
-- ) VALUES (
--   $1, $2
-- )
-- RETURNING *;

-- name: DeleteAuthor :exec
-- DELETE FROM authors
-- WHERE id = $1;

-- Schema:
-- CREATE TABLE authors (
--   id   BIGSERIAL PRIMARY KEY,
--   name text      NOT NULL,
--   bio  text
-- );
