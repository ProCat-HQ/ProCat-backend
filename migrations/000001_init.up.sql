CREATE TABLE IF NOT EXISTS users
(
    id                    SERIAL PRIMARY KEY,
    fullname              VARCHAR            NOT NULL,
    email                 VARCHAR(255) UNIQUE,
    phone_number          VARCHAR(20) UNIQUE NOT NULL,
    identification_number VARCHAR(20) UNIQUE,
    password_hash         VARCHAR            NOT NULL,
    is_confirmed          BOOLEAN            NOT NULL DEFAULT FALSE,
    role                  VARCHAR(30)        NOT NULL DEFAULT 'user',
    created_at            TIMESTAMP                   DEFAULT now()
);

CREATE TABLE IF NOT EXISTS refresh_sessions
(
    id            SERIAL PRIMARY KEY,
    refresh_token VARCHAR NOT NULL,
    fingerprint   VARCHAR NOT NULL,
    user_id       INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS deliverymen
(
    id                  SERIAL PRIMARY KEY,
    car_capacity        VARCHAR(255),
    working_hours_start TIME           NOT NULL,
    working_hours_end   TIME           NOT NULL,
    car_id              VARCHAR(30),
    user_id             INTEGER UNIQUE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS verifications
(
    id        SERIAL PRIMARY KEY,
    code      VARCHAR   NOT NULL,
    type      VARCHAR   NOT NULL,
    value     VARCHAR   NOT NULL,
    life_time TIMESTAMP NOT NULL DEFAULT now() + interval '5' MINUTE,
    user_id   INTEGER   NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orders
(
    id                  SERIAL PRIMARY KEY,
    status              VARCHAR(40) NOT NULL,
    total_price         INTEGER     NOT NULL,
    deposit             INTEGER,
    rental_period_start TIMESTAMP   NOT NULL,
    rental_period_end   TIMESTAMP   NOT NULL,
    address             VARCHAR     NOT NULL,
    latitude            VARCHAR,
    longitude           VARCHAR,
    company_name        VARCHAR(255),
    created_at          TIMESTAMP DEFAULT now(),
    user_id             INTEGER     NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS deliveries
(
    id             SERIAL PRIMARY KEY,
    time_start     TIMESTAMP   NOT NULL,
    time_end       TIMESTAMP   NOT NULL,
    method         VARCHAR(50) NOT NULL,
    order_id       INTEGER     NOT NULL,
    deliveryman_id INTEGER,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    FOREIGN KEY (deliveryman_id) REFERENCES deliverymen (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS routes
(
    id             SERIAL PRIMARY KEY,
    deliveryman_id INTEGER,
    FOREIGN KEY (deliveryman_id) REFERENCES deliverymen (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS coordinates
(
    id              SERIAL PRIMARY KEY,
    sequence_number INTEGER NOT NULL,
    delivery_id        INTEGER NOT NULL,
    route_id        INTEGER NOT NULL,
    FOREIGN KEY (delivery_id) REFERENCES deliveries (id) ON DELETE SET NULL,
    FOREIGN KEY (route_id) REFERENCES routes (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payments
(
    id         SERIAL PRIMARY KEY,
    paid       INTEGER     NOT NULL DEFAULT 0,
    method     VARCHAR(50),
    price      INTEGER     NOT NULL,
    created_at TIMESTAMP            DEFAULT now(),
    order_id   INTEGER     NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS categories
(
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(255),
    parent_id INTEGER,
    FOREIGN KEY (parent_id) REFERENCES categories (id) ON DELETE CASCADE -- OR MAYBE LEAVE IT JUST LIKE INTEGER
);

CREATE TABLE IF NOT EXISTS items
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    description   VARCHAR,
    price         INTEGER      NOT NULL,
    price_deposit INTEGER      NOT NULL,
    is_in_stock   BOOLEAN      NOT NULL DEFAULT FALSE,
    category_id   INTEGER,
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET NULL
--     similar_to INTEGER[],
--     FOREIGN KEY (EACH ELEMENT OF similar_to) REFERENCES items
);

CREATE TABLE IF NOT EXISTS stores
(
    id                  SERIAL PRIMARY KEY,
    name                VARCHAR NOT NULL,
    address             VARCHAR NOT NULL,
    latitude            VARCHAR,
    longitude           VARCHAR,
    working_hours_start TIME    NOT NULL,
    working_hours_end   TIME    NOT NULL
);

CREATE TABLE IF NOT EXISTS item_stores
(
    id              SERIAL PRIMARY KEY,
    in_stock_number INTEGER NOT NULL,
    store_id        INTEGER,
    item_id         INTEGER NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE,
    FOREIGN KEY (store_id) REFERENCES stores (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS item_images
(
    id      SERIAL PRIMARY KEY,
    image   VARCHAR NOT NULL,
    item_id INTEGER NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orders_items
(
    id           SERIAL PRIMARY KEY,
    items_number INTEGER NOT NULL DEFAULT 1,
    order_id     INTEGER NOT NULL,
    item_id      INTEGER, -- ОБГОВОРИТЬ (ЧТО ЕСЛИ ТОВАР УДАЛЯЕТСЯ, А ЗАКАЗ ОСТАЁТСЯ)
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS infos
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description VARCHAR      NOT NULL,
    item_id     INTEGER      NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS carts
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS carts_items
(
    id           SERIAL PRIMARY KEY,
    items_number INTEGER NOT NULL DEFAULT 1,
    cart_id      INTEGER NOT NULL,
    item_id      INTEGER,
    FOREIGN KEY (cart_id) REFERENCES carts (id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subscriptions
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subscriptions_items
(
    id              SERIAL PRIMARY KEY,
    subscription_id INTEGER NOT NULL,
    item_id         INTEGER,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions (id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notifications
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR      NOT NULL,
    is_viewed   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP             DEFAULT now(),
    user_id     INTEGER      NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chats
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    is_solved      BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMP             DEFAULT now(),
    first_user_id  INTEGER,
    second_user_id INTEGER,
    order_id       INTEGER,
    FOREIGN KEY (first_user_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (second_user_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages
(
    id         SERIAL PRIMARY KEY,
    text       VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    user_id    INTEGER NOT NULL, -- ЧТО ДЕЛАТЬ ДЛЯ АНОНИМНОГО ПОЛЬЗОВАТЕЛЯ
    chat_id    INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (chat_id) REFERENCES chats (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS message_images
(
    id         SERIAL PRIMARY KEY,
    image      VARCHAR NOT NULL,
    message_id INTEGER NOT NULL,
    FOREIGN KEY (message_id) REFERENCES messages (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_index_carts_items ON carts_items (cart_id, item_id);
CREATE UNIQUE INDEX IF NOT EXISTS unique_index_subscriptions_items ON subscriptions_items (subscription_id, item_id);
CREATE UNIQUE INDEX IF NOT EXISTS unique_index_orders_items ON orders_items (order_id, item_id);
CREATE UNIQUE INDEX IF NOT EXISTS unique_index_item_stores ON item_stores (store_id, item_id);