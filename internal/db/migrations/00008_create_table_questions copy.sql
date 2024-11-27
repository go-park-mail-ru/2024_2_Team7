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

CREATE TABLE IF NOT EXISTS ANSWERS(
                            id SERIAL PRIMARY KEY,
                            question_id INT NOT NULL,
                            user_id INT NOT NULL,
                            answer INT,
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            FOREIGN KEY (question_id) REFERENCES QUESTION (id) ON DELETE CASCADE,
                            FOREIGN KEY (user_id) REFERENCES "USER" (id) ON DELETE CASCADE,
                            CONSTRAINT unique_interviewed UNIQUE (question_id, user_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ANSWERS;
DROP TABLE IF EXISTS QUESTION;
DROP TABLE IF EXISTS TEST;
-- +goose StatementEnd
