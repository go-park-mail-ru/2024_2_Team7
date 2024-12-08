# ДЗ1 СУБД
## Безопасность сервера СУБД

**Создание пользователя с ограниченными правами**:  
Для создания пользователя используется скрипт с правами только на чтение данных из нужных таблиц:

   ```sql
   CREATE ROLE ${POSTGRES_USER} WITH LOGIN PASSWORD '${POSTGRES_PASSWORD}';  
GRANT CONNECT ON DATABASE ${POSTGRES_DB} TO ${POSTGRES_USER};  
GRANT USAGE ON SCHEMA public TO ${POSTGRES_USER};  
GRANT SELECT ON ALL TABLES IN SCHEMA public TO ${POSTGRES_USER};   
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO ${POSTGRES_USER};  
GRANT USAGE ON SCHEMA public TO ${POSTGRES_USER};
   ```
Скрипт создается в репозитории в директории `/internal/db/` для удобства управления правами на БД

## Защита от SQL инъекций
**Prepared Statements**:  
Используется подготовленный запрос с параметрами, что защищает от SQL инъекций. Пример запроса с использованием плейсхолдеров в Go:
```go
const createEventQuery = `
INSERT INTO event (title, description, event_start, event_finish, location, capacity, user_id, category_id, lat, lon)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id`

```
## Пул соединений и параметры соединений
**Connection Pooling**:
Выбор значения `max_connections` аргументирован с учетом малого объема оперативной памяти (3 ГБ) на сервере. Излишнее количество соединений может привести к избыточному потреблению ресурсов, поэтому было выбрано небольшое значение с учетом нагрузки.
`max_connections = 100`

**Настройка listen_addresses**:
Значение `listen_addresses` настроено на `localhost`, что минимизирует риски внешнего доступа и повышает безопасность. Так как все контейнеры приложения работают на одной физической вм.

## Настройка параметров сервера и клиента
**Таймауты**:  
Параметры `statement_timeout` и `lock_timeout` настроены с учетом бизнес требований. Для предотвращения DoS атак установлены ограничения:
```conf
statement_timeout = 30000  # 30 секунд
lock_timeout = 5000        # 5 секунд
```
- `statement_timeout` установлен в 30 секунд, что является разумным значением для выполнения запросов. Это значение достаточно для большинства операций, но если запросы будут длиться дольше, это может сигнализировать о проблемах с производительностью. Ограничение в 30 секунд снижает вероятность зависания приложения из-за долгих запросов.
- `lock_timeout` установлен в 5 секунд, что помогает избежать блокировок и замедления системы при длительных ожиданиях на получение блокировки. Эти значения обеспечивают баланс между производительностью и стабильностью.

**pg_stat_statements**:  
Для логирования медленных запросов и их анализа настроены параметры:
```conf
shared_preload_libraries = 'pg_stat_statements, auto_explain'
log_statement = 'none'
log_min_duration_statement = 1000  # Логировать запросы, если время выполнения больше 1 секунды
```
Включено логирование медленных запросов для анализа и предотвращения проблем с производительностью. Используется PGBadger для анализа логов:
```conf
logging_collector = on  
log_directory = '/var/log/postgresql'  
log_filename = 'postgresql-%a.log'  
log_rotation_age = 1d  
log_rotation_size = 100MB  
log_min_duration_statement = 1000  # Время выполнения запроса более 1 секунды
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h'
```
