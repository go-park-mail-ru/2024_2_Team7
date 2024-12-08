-- +goose Up
-- +goose StatementBegin
CREATE TABLE NOTIFICATION (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    event_id INT NOT NULL,
    notify_at TIMESTAMP NOT NULL, 
    is_sent BOOLEAN DEFAULT FALSE, 
    message TEXT NOT NULL, 
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS NOTIFICATION;
-- +goose StatementEnd
