-- name: GetUserByEmail :one
SELECT id
      ,email
      ,"password"
      ,first_name
      ,last_name
      ,created_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (email, "password", first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING *;