-- name: CreatNote :one
INSERT INTO notes (user_id, title, content)
VALUES ( $1, $2, $3 )
RETURNING *;

-- name: ListNotesByUserID :many
SELECT * FROM notes
WHERE user_id = $1;

-- name: GetNoteByUserID :one
SELECT * FROM notes
WHERE id = $1 AND user_id = $2;

-- name: UpdateNote :one
UPDATE notes
SET title = $2, content = $3, updated_at = now()
WHERE id = $1 AND user_id = $4
RETURNING *;


-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = $1 AND user_id = $2;
