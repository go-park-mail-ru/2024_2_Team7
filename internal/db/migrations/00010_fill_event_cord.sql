-- +goose Up
-- +goose StatementBegin
DO $$
    DECLARE
        event_id INT;
        generated_lat DOUBLE PRECISION;
        generated_lon DOUBLE PRECISION;
    BEGIN
        -- Перебор всех событий с id от 1 до 32
        FOR event_id IN 1..32 LOOP
                -- Генерация случайных координат в пределах радиуса 100 км от центра Москвы
                generated_lat := 55.7558 + (random() - 0.5) * 0.9;  -- Москва: широта 55.7558, разброс ±0.45
                generated_lon := 37.6173 + (random() - 0.5) * 1.2;  -- Москва: долгота 37.6173, разброс ±0.6

                -- Обновление событий с вычисленными координатами
                UPDATE EVENT
                SET
                    lat = generated_lat,
                    lon = generated_lon
                WHERE id = event_id;
            END LOOP;
    END $$;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
