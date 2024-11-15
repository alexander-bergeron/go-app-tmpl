-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :exec
INSERT INTO users (username, email, first_name, last_name)
VALUES ($1, $2, $3, $4);
