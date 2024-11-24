-- +goose Up
-- +goose StatementBegin
INSERT INTO CATEGORY (name) VALUES
('Музыка'),
('Спорт'),
('Выставки'),
('Образование'),
('Театр'),
('Кино'),
('Еда'),
('Для детей');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM CATEGORY;
-- +goose StatementEnd
