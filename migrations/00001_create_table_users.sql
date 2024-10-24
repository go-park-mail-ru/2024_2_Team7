-- +goose Up
-- +goose StatementBegin
CREATE TABLE "USER" (
                        id SERIAL PRIMARY KEY,
                        username TEXT NOT NULL,
                        email TEXT UNIQUE NOT NULL,
                        password_hash TEXT NOT NULL,
                        URL_to_avatar TEXT,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "USER";
-- +goose StatementEnd
