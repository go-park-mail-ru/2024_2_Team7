-- переменные берем из .env, которые передаем в docker-compose.yml
-- Создание пользователя с ограниченными правами
CREATE ROLE ${POSTGRES_USER} WITH LOGIN PASSWORD '${POSTGRES_PASSWORD}';
GRANT CONNECT ON DATABASE ${POSTGRES_DB} TO ${POSTGRES_USER};

-- Ограничение прав на чтение данных (SELECT)
GRANT USAGE ON SCHEMA public TO ${POSTGRES_USER};
GRANT SELECT ON ALL TABLES IN SCHEMA public TO ${POSTGRES_USER};

-- Обеспечение прав для будущих таблиц
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO ${POSTGRES_USER};

-- Дать доступ на использование схемы
GRANT USAGE ON SCHEMA public TO ${POSTGRES_USER};
