-- name: CreateUser :one
INSERT INTO users (id, create_at, updated_at, display_name, email_address, password, api_key)
VALUES ($1, $2, $3, $4, $5, $6,
    encode(sha256(random()::text::bytea), 'hex')
)
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;

-- name: LogIn :one
SELECT * FROM users u WHERE u.email_address = $1 and u.password = $2;