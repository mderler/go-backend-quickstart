-- +goose Up
-- +goose StatementBegin
ALTER TABLE todo RENAME COLUMN user_id TO creator_id;

CREATE TABLE todo_user (
    todo_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (todo_id) REFERENCES todo(id),
    FOREIGN KEY (user_id) REFERENCES "user"(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE todo RENAME COLUMN creator_id TO user_id;
DROP TABLE todo_user;
-- +goose StatementEnd
