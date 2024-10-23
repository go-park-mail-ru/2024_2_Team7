CREATE TABLE "USER" (
                        id SERIAL PRIMARY KEY,
                        username TEXT NOT NULL,
                        email TEXT UNIQUE NOT NULL,
                        password_hash TEXT NOT NULL,
                        URL_to_avatar TEXT,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE CATEGORY (
                          id SERIAL PRIMARY KEY,
                          name TEXT NOT NULL
);

CREATE TABLE EVENT (
                       id SERIAL PRIMARY KEY,
                       title TEXT NOT NULL,
                       description TEXT,
                       event_start DATE NOT NULL,
                       event_finish DATE NOT NULL,
                       location TEXT,
                       capacity INT,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_DATE,
                       user_id INT NOT NULL,
                       category_id INT,
                       FOREIGN KEY (user_id) REFERENCES "USER" (id) ON DELETE CASCADE,
                       FOREIGN KEY (category_id) REFERENCES CATEGORY (id) ON DELETE CASCADE
);

CREATE TABLE MEDIA_URL (
                           id SERIAL PRIMARY KEY,
                           url TEXT NOT NULL,
                           event_id INT NOT NULL,
                           FOREIGN KEY (event_id) REFERENCES EVENT (id) ON DELETE CASCADE
);

CREATE TABLE TAG (
                     id SERIAL PRIMARY KEY,
                     name TEXT NOT NULL,
                     created_at DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE TABLE EVENT_TAG (
                           id SERIAL PRIMARY KEY,
                           event_id INT NOT NULL,
                           tag_id INT NOT NULL,
                           FOREIGN KEY (event_id) REFERENCES EVENT (id) ON DELETE CASCADE,
                           FOREIGN KEY (tag_id) REFERENCES TAG (id) ON DELETE CASCADE
);

CREATE TABLE TICKET (
                        id SERIAL PRIMARY KEY,
                        ticket_buy_date DATE NOT NULL DEFAULT CURRENT_DATE,
                        type TEXT NOT NULL,
                        price DECIMAL(10, 2) NOT NULL,
                        quantity INT NOT NULL,
                        event_id INT NOT NULL,
                        user_id INT NOT NULL,
                        FOREIGN KEY (event_id) REFERENCES EVENT (id) ON DELETE CASCADE,
                        FOREIGN KEY (user_id) REFERENCES "USER" (id) ON DELETE CASCADE
);

CREATE TABLE ATTENDANCE (
                            user_id INT NOT NULL,
                            event_id INT NOT NULL,
                            attended_at DATE NOT NULL DEFAULT CURRENT_DATE,
                            PRIMARY KEY (user_id, event_id),
                            FOREIGN KEY (user_id) REFERENCES "USER" (id) ON DELETE CASCADE,
                            FOREIGN KEY (event_id) REFERENCES EVENT (id) ON DELETE CASCADE
);
