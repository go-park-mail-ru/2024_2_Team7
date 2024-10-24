-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO CATEGORY (name) VALUES
('music'),
('sport'),
('exhibitions'),
('education'),
('theater'),
('cinema'),
('food'),
('children');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM CATEGORY;
-- +goose StatementEnd
