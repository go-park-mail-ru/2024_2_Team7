-- +goose Up
-- +goose StatementBegin
CREATE TABLE FAVORITE_EVENT (
                            id SERIAL PRIMARY KEY,
                            user_id INT NOT NULL,
                            event_id INT NOT NULL,
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            FOREIGN KEY (user_id) REFERENCES "USER" (id) ON DELETE CASCADE,
                            FOREIGN KEY (event_id) REFERENCES EVENT (id) ON DELETE CASCADE,
                            CONSTRAINT unique_favorites UNIQUE (user_id, event_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS FAVORITE_EVENT;
-- +goose StatementEnd
