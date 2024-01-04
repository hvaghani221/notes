-- name: CreatNote :one
INSERT INTO notes (user_id, title, content)
VALUES ( $1, $2, $3 )
RETURNING *;

-- name: ListNotesByUserID :many
SELECT n.id, n.user_id, n.title, n.content, n.created_at, n.updated_at
FROM notes n
LEFT JOIN shared_notes sn ON n.id = sn.note_id
WHERE n.user_id = $1 OR sn.shared_with_user_id = $1
ORDER BY n.created_at DESC;

-- name: GetNoteByUserID :one
SELECT n.id, n.user_id, n.title, n.content, n.created_at, n.updated_at
FROM notes n
LEFT JOIN shared_notes sn ON n.id = sn.note_id
WHERE (n.id = $1) AND (n.user_id = $2 OR sn.shared_with_user_id = $2);

-- name: UpdateNote :one
UPDATE notes
SET title = $2, content = $3, updated_at = now()
WHERE id = $1 AND user_id = $4
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = $1 AND user_id = $2;

-- name: ShareNote :one
INSERT INTO shared_notes (note_id, shared_with_user_id)
SELECT @noteID, users.id
FROM users
WHERE email = @sharedWithEmail AND EXISTS (SELECT 1 FROM notes WHERE id = @noteID AND user_id = @userID)
RETURNING *;
