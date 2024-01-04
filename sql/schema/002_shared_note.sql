-- +goose Up
CREATE TABLE shared_notes (
    note_id INT REFERENCES notes(id) ON DELETE CASCADE,
    shared_with_user_id INT REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (note_id, shared_with_user_id)
);

-- +goose Down
DROP TABLE shared_notes;
