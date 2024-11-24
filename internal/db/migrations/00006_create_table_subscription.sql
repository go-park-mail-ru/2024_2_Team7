-- +goose Up
-- +goose StatementBegin
CREATE TABLE SUBSCRIPTION (
                            id SERIAL PRIMARY KEY,
                            subscriber_id INT NOT NULL,
                            subscribed_id INT NOT NULL,
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            FOREIGN KEY (subscriber_id) REFERENCES "USER" (id) ON DELETE CASCADE,
                            FOREIGN KEY (subscribed_id) REFERENCES "USER" (id) ON DELETE CASCADE,
                            CONSTRAINT unique_subscription UNIQUE (subscriber_id, subscribed_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS SUBSCRIPTION;
-- +goose StatementEnd
