-- +goose Up
-- +goose StatementBegin
INSERT INTO test (title)
VALUES ('test 1');

INSERT INTO question (test_id, question)
VALUES ('1', 'who?'), ('1', 'why?');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
