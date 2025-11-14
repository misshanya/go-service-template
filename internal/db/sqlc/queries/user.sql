-- name: CreateUser :one
INSERT INTO users (username, age)
VALUES (
        @username::text, @age::int
)
RETURNING *;
