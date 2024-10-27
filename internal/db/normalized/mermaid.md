```mermaid
erDiagram
    USER {
        int id PK "Primary Key"
        TEXT name 
        TEXT email 
        TEXT password 
        date created_at 
        TEXT URL_to_avatar }

    EVENT {
        int id PK "Primary Key"
        TEXT title 
        TEXT description 
        date event_start
        date event_finish
        TEXT location 
        int capacity 
        date created_at
        int user_id FK "Foreign Key to USER"
        int category_id FK "Foreign Key to CATEGORY"
    }

    MEDIA_URL {
        int id PK
        TEXT url
        int event_id FK
    }

    TAG {
        int id PK "Primary Key"
        TEXT name 
        date created_at
    }

    EVENT_TAG {
        int id PK
        int event_id FK "Foreign Key to EVENT"
        int tag_id FK "Foreign Key to TAG"
    }

    CATEGORY {
        int id PK "Primary Key"
        TEXT name }

    TICKET {
        int id PK "Primary Key"
        date ticket_buy_date 
        TEXT type 
        decimal price 
        int quantity 
        int event_id FK "Foreign Key to EVENT"
        int user_id FK "Foreign key to USER"
    }

    ATTENDANCE {
        int user_id FK "Foreign Key to USER"
        int event_id FK "Foreign Key to EVENT"
        date attended_at
    }

    USER ||--o{ EVENT : creates 
    USER ||--o{ TICKET : has

    EVENT ||--o{ ATTENDANCE : attended_by
    EVENT ||--o| CATEGORY : belongs_to
    EVENT ||--o{ EVENT_TAG : has 

    TAG ||--o{ EVENT_TAG : includes

    EVENT ||--o{ MEDIA_URL : includes