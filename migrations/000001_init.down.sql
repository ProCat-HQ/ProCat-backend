DROP INDEX IF EXISTS unique_index_carts_items;
DROP INDEX IF EXISTS unique_index_subscriptions_items;
DROP INDEX IF EXISTS unique_index_orders_items;
DROP INDEX IF EXISTS unique_index_item_stores;

DROP TABLE IF EXISTS message_images;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;

DROP TABLE IF EXISTS notifications;

DROP TABLE IF EXISTS subscriptions_items;
DROP TABLE IF EXISTS subscriptions;

DROP TABLE IF EXISTS carts_items;
DROP TABLE IF EXISTS carts;

DROP TABLE IF EXISTS infos;
DROP TABLE IF EXISTS orders_items;
DROP TABLE IF EXISTS item_images;
DROP TABLE IF EXISTS item_stores;
DROP TABLE IF EXISTS stores;

DROP TABLE IF EXISTS items;

DROP TABLE IF EXISTS categories;

DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS coordinates;
DROP TABLE IF EXISTS routes;

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS verifications;
DROP TABLE IF EXISTS deliverymen;
DROP TABLE IF EXISTS refresh_sessions;
DROP TABLE IF EXISTS users;
