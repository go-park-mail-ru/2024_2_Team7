-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS TEST(
                            id SERIAL PRIMARY KEY,
                            title text,                           
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS QUESTION(
                            id SERIAL PRIMARY KEY,
                            test_id INT NOT NULL,
                            question text,             
                            FOREIGN KEY (test_id) REFERENCES TEST (id) ON DELETE CASCADE,              
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS INTERVIEWED(
                            id SERIAL PRIMARY KEY,
                            test_id INT NOT NULL,
                            user_id INT NOT NULL,
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            FOREIGN KEY (test_id) REFERENCES TEST (id) ON DELETE CASCADE,
                            FOREIGN KEY (user_id) REFERENCES "USER" (id) ON DELETE CASCADE,
                            CONSTRAINT unique_interviewed UNIQUE (test_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS TEST;
DROP TABLE IF EXISTS QUESTION;
DROP TABLE IF EXISTS INTERVIEWED;
-- +goose StatementEnd
