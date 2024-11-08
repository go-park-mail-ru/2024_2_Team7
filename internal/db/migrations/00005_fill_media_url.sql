-- +goose Up
-- +goose StatementBegin
INSERT INTO MEDIA_URL (url, event_id)
VALUES
    ('static/images/1.png', 1),
    ('static/images/2.png', 2),
    ('static/images/3.png', 3),
    ('static/images/4.png', 4),
    ('static/images/5.png', 5),
    ('static/images/6.png', 6),
    ('static/images/7.png', 7),
    ('static/images/8.png', 8),
    ('static/images/9.png', 9),
    ('static/images/10.png', 10),
    ('static/images/11.png', 11),
    ('static/images/12.png', 12),
    ('static/images/13.png', 13),
    ('static/images/14.png', 14),
    ('static/images/15.png', 15),
    ('static/images/16.png', 16),
    ('static/images/17.png', 17),
    ('static/images/18.png', 18),
    ('static/images/19.png', 19),
    ('static/images/20.png', 20),
    ('static/images/21.png', 21),
    ('static/images/22.png', 22),
    ('static/images/23.png', 23),
    ('static/images/24.png', 24),
    ('static/images/25.png', 25),
    ('static/images/26.png', 26),
    ('static/images/27.png', 27),
    ('static/images/28.png', 28),
    ('static/images/29.png', 29),
    ('static/images/30.png', 30),
    ('static/images/31.png', 31),
    ('static/images/32.png', 32);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
