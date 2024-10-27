ДЗ1 СУБД

диаграмму удобнее вставлять *[сюда](https://mermaid.live/)*
```
erDiagram
    USER {
        int id PK NOT NULL
        TEXT name NOT NULL
        TEXT email UNIQUE NOT NULL
        TEXT password NOT NULL
        date created_at DEFAULT CURRENT_DATE
        TEXT URL_to_avatar 
    }

    EVENT {
        int id PK NOT NULL
        TEXT title NOT NULL
        TEXT description 
        date event_start NOT NULL
        date event_finish NOT NULL
        TEXT location 
        int capacity DEFAULT 0
        date created_at DEFAULT CURRENT_DATE
        int user_id FK NOT NULL ON DELETE CASCADE ON UPDATE CASCADE
        int category_id FK NOT NULL ON DELETE CASCADE ON UPDATE CASCADE
    }

    MEDIA_URL {
        int id PK NOT NULL
        TEXT url NOT NULL
        int event_id FK NOT NULL ON DELETE CASCADE ON UPDATE CASCADE
    }

    TAG {
        int id PK NOT NULL
        TEXT name UNIQUE NOT NULL
        date created_at DEFAULT CURRENT_DATE
    }

    EVENT_TAG {
        int id PK NOT NULL
        int event_id FK NOT NULL ON DELETE CASCADE ON UPDATE CASCADE
        int tag_id FK NOT NULL ON DELETE CASCADE ON UPDATE CASCADE
        UNIQUE (event_id, tag_id)
    }

    CATEGORY {
        int id PK NOT NULL
        TEXT name UNIQUE NOT NULL
    }

    TICKET {
        int id PK NOT NULL
        date ticket_buy_date NOT NULL
        TEXT type NOT NULL
        decimal price NOT NULL CHECK(price >= 0)
        int quantity DEFAULT 1 CHECK(quantity > 0)
        int event_id FK NOT NULL ON DELETE CASCADE ON UPDATE CASCADE
        int user_id FK NOT NULL ON DELETE CASCADE ON UPDATE CASCADE
    }

    ATTENDANCE {
        int user_id FK NOT NULL ON DELETE CASCADE ON UPDATE CASCADE
        int event_id FK NOT NULL ON DELETE CASCADE ON UPDATE CASCADE
        date attended_at NOT NULL
        UNIQUE (user_id, event_id)
    }

```
1. Первая нормальная форма (1NF) Для того чтобы таблица соответствовала 1NF, все атрибуты должны содержать только атомарные (неделимые) значения, и каждая запись должна быть уникальной.
   Во всех таблицах соблюдены правила 1NF, потому что:
    - Каждое поле хранит атомарные (неделимые) данные.
    - В таблицах отсутствуют повторяющиеся группы данных и списки значений в одном поле.
    - У каждой таблицы есть первичный ключ, обеспечивающий уникальность записей.
    - Внешние ключи используются для создания связей между таблицами, исключая дублирование данных.
2. Вторая нормальная форма (2NF) Для того чтобы таблица соответствовала 2NF, она должна соответствовать 1NF и все неключевые атрибуты (поля, не входящие в состав первичного ключа) должны зависеть от всего составного ключа, а не только от его части.
   Схема соответствует 2NF так как:
    - Все таблицы находятся в 1NF.
    - В таблицах с одинарными ключами нет частичных зависимостей.
    - В таблице с составным ключом (`EVENT_TAG`) нет полей, зависящих от части ключа, что исключает частичные зависимости.
3. Третья нормальная форма (3NF) Для того чтобы таблица соответствовала 3NF, она должна соответствовать 2NF,и каждый неключевой атрибут должен зависеть только от первичного ключа и ни от каких других неключевых атрибутов.
   Схема соответствует 3NF так как:
    - Все таблицы приведены к 2NF.
    - Во всех таблицах отсутствуют транзитивные зависимости: неключевые атрибуты зависят исключительно от первичного ключа и не зависят друг от друга.
	