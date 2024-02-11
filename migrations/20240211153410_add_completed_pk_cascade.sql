-- +goose Up
-- +goose StatementBegin
ALTER TABLE todo ADD completed BOOLEAN DEFAULT FALSE;

ALTER TABLE todo_user
DROP CONSTRAINT todo_user_todo_id_fkey,
DROP CONSTRAINT todo_user_user_id_fkey,
ADD CONSTRAINT todo_user_pk PRIMARY KEY (todo_id, user_id),
ADD CONSTRAINT todo_user_todo_id_fkey FOREIGN KEY (todo_id) REFERENCES todo(id) ON DELETE CASCADE,
ADD CONSTRAINT todo_user_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE todo DROP COLUMN completed;

ALTER TABLE todo_user
DROP CONSTRAINT todo_user_pk,
DROP CONSTRAINT todo_user_todo_id_fkey,
DROP CONSTRAINT todo_user_user_id_fkey,
ADD CONSTRAINT todo_user_todo_id_fkey FOREIGN KEY (todo_id) REFERENCES todo(id),
ADD CONSTRAINT todo_user_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user"(id);
-- +goose StatementEnd
