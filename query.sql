-- name: ListPosts :many
SELECT id, title, subtitle, content, created_at, updated_at FROM Posts;

-- name: GetPost :one
SELECT id, title, subtitle, content, created_at, updated_at FROM Posts WHERE id = $1;

-- name: InsertPost :exec
INSERT INTO Posts (id, title, subtitle, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdatePost :exec
UPDATE Posts SET title = $1, subtitle = $2, content = $3, updated_at = $4 WHERE id = $5;

-- name: DeletePost :exec
DELETE FROM Posts WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, password FROM Users WHERE email = $1;

-- name: InsertUser :exec
INSERT INTO Users (id, name, email, password) VALUES ($1, $2, $3, $4);
