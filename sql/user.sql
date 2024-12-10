-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (username, email, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET username = $2, email = $3, first_name = $4, last_name = $5, version = version + 1
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE user_id = $1
RETURNING *;

