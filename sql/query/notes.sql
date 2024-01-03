-- name: CreatNote :one
INSERT INTO notes (user_id, title, content)
VALUES ( $1, $2, $3 )
RETURNING *;

-- name: GetNotesByUserID :many
SELECT * FROM notes
WHERE user_id = $1;
