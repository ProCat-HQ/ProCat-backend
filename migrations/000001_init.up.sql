CREATE TABLE users
(
    id                    SERIAL PRIMARY KEY,
    fullName              VARCHAR            NOT NULL,
    email                 VARCHAR(255) UNIQUE,
    phone_number          VARCHAR(20) UNIQUE NOT NULL,
    identification_number VARCHAR(20) UNIQUE,
    password_hash         VARCHAR            NOT NULL,
    is_confirmed          BOOLEAN            NOT NULL DEFAULT FALSE,
    role                  VARCHAR(30)        NOT NULL DEFAULT 'user',
    created_at            TIMESTAMP          NOT NULL
);

CREATE TABLE delivery_men
(
    id                  SERIAL PRIMARY KEY,
    car_capacity        VARCHAR(255),
    working_hours_start TIME    NOT NULL,
    working_hours_end   TIME    NOT NULL,
    car_id              VARCHAR(30),
    user_id             INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE orders
(
    id                  SERIAL PRIMARY KEY,
    status              VARCHAR(40) NOT NULL,
    total_price         INTEGER     NOT NULL,
    rental_period_start TIMESTAMP   NOT NULL,
    rental_period_end   TIMESTAMP,                                 -- NOT NULL?
    address             VARCHAR     NOT NULL,
    company_name        VARCHAR(255),
--     contract            TEXT,
    user_id             INTEGER     NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL -- OR DELETE?
);

CREATE TABLE deliveries
(
    id              SERIAL PRIMARY KEY,
    time            TIMESTAMP   NOT NULL,
    method          VARCHAR(50) NOT NULL,
    order_id        INTEGER     NOT NULL,
    delivery_man_id INTEGER,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    FOREIGN KEY (delivery_man_id) REFERENCES delivery_men (id) ON DELETE SET NULL
);

CREATE TABLE payments
(
    id       SERIAL PRIMARY KEY,
    is_paid  BOOLEAN     NOT NULL DEFAULT FALSE,
    method   VARCHAR(50) NOT NULL,
    price    INTEGER     NOT NULL,
    order_id INTEGER     NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

CREATE TABLE categories
(
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(255),
    parent_id INTEGER,
    FOREIGN KEY (parent_id) REFERENCES categories (id) ON DELETE CASCADE -- OR MAYBE LEAVE IT JUST LIKE INTEGER
);

CREATE TABLE items
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description VARCHAR,
    price       INTEGER      NOT NULL,
    category_id INTEGER,
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET NULL
--     similar_to INTEGER[],
--     FOREIGN KEY (EACH ELEMENT OF similar_to) REFERENCES items
);

CREATE TABLE item_statuses
(
    id              SERIAL PRIMARY KEY,
    is_in_stock     BOOLEAN NOT NULL,
    in_stock_number INTEGER NOT NULL,
    item_id         INTEGER NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
);

CREATE TABLE item_images
(
    id      SERIAL PRIMARY KEY,
    image   VARCHAR NOT NULL,
    item_id INTEGER NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
);

CREATE TABLE orders_items
(
    id           SERIAL PRIMARY KEY,
    items_number INTEGER NOT NULL DEFAULT 1,
    order_id     INTEGER NOT NULL,
    item_id      INTEGER, -- ОБГОВОРИТЬ (ЧТО ЕСЛИ ТОВАР УДАЛЯЕТСЯ, А ЗАКАЗ ОСТАЁТСЯ)
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE SET NULL
);

CREATE TABLE infos
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description VARCHAR      NOT NULL,
    item_id     INTEGER      NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
);

CREATE TABLE carts
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE carts_items
(
    id           SERIAL PRIMARY KEY,
    items_number INTEGER NOT NULL DEFAULT 1,
    cart_id      INTEGER NOT NULL,
    item_id      INTEGER,
    FOREIGN KEY (cart_id) REFERENCES carts (id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE SET NULL
);

CREATE TABLE subscriptions
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE subscriptions_items
(
    id              SERIAL PRIMARY KEY,
    subscription_id INTEGER NOT NULL,
    item_id         INTEGER,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions (id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE SET NULL
);

CREATE TABLE notifications
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR      NOT NULL,
    is_viewed   BOOLEAN      NOT NULL DEFAULT FALSE,
    user_id     INTEGER      NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE chats
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    is_solved      BOOLEAN      NOT NULL DEFAULT FALSE,
    first_user_id  INTEGER,
    second_user_id INTEGER,
    order_id       INTEGER,
    FOREIGN KEY (first_user_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (second_user_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

CREATE TABLE messages
(
    id      SERIAL PRIMARY KEY,
    text    VARCHAR NOT NULL,
    user_id INTEGER NOT NULL, -- ЧТО ДЕЛАТЬ ДЛЯ АНОНИМНОГО ПОЛЬЗОВАТЕЛЯ
    chat_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (chat_id) REFERENCES chats (id) ON DELETE CASCADE
);

CREATE TABLE message_images
(
    id         SERIAL PRIMARY KEY,
    image      VARCHAR NOT NULL,
    message_id INTEGER NOT NULL,
    FOREIGN KEY (message_id) REFERENCES messages (id) ON DELETE CASCADE
);