INSERT INTO users (id, fullname, email, phone_number, identification_number, password_hash, role)
VALUES (100, 'Админ Админович', 'admin@gmail.com', '+79998887766', '123',
        '6e77446d383649393b4c653428734456781324630821b58653d4a6e43aa7d28afc552043', 'admin'),

       (101, 'Доставщик Первый Доставщикович', 'deliveryman1@gmail.com', '+71998887766', '234',
        '6e77446d383649393b4c653428734456c499bc1accc120a80337058194742db0b3c6e13a', 'deliveryman'),

       (102, 'Доставщик Второй Доставщикович', 'deliveryman2@gmail.com', '+71298887766', '2344',
        '6e77446d383649393b4c653428734456d46fbd98aa243d1d58c7132ac8e7a29ca90f7e31', 'deliveryman'),

       (103, 'User Usernamovich', 'user@gmail.com', '+72998887766', '345',
        '6e77446d383649393b4c6534287344564f61d67813c5a1818322c1cf96890b8b7934f97f', 'user');

INSERT INTO carts (user_id)
VALUES (100),
       (101),
       (102),
       (103);

INSERT INTO deliverymen (car_capacity, working_hours_start, working_hours_end, car_id, user_id)
VALUES ('big', '10:00', '16:00', 'A123BC', 101),
       ('medium', '12:00', '20:00', 'X777XX', 102);

INSERT INTO categories (id, name, parent_id)
VALUES (1, 'cat1', 2),
       (2, 'cat2', 1);

INSERT INTO items (name, description, price, price_deposit, is_in_stock, category_id)
VALUES ('Молоток', 'Хороший молоток, крепкий', 2000, 20000, true, 2),
       ('Пила', 'Хорошо пилит, мощно', 3000, 30000, false, null),
       ('Плоскогубцы', 'Хорошо сжимает, крепкая хватка', 4000, 40000, true, 1);

INSERT INTO infos (name, description, item_id)
VALUES ('ВЕС', 'большой', 1),
       ('gabariti', 'like a closet', 1),
       ('eshcho cho-ta', 'property', 1),
       ('eshcho cho-ta', 'property', 2);

INSERT INTO item_images (image, item_id)
VALUES ('hammer.jpg', 1),
       ('hammer.jpg', 1),
       ('set.jpg', 3);

INSERT INTO stores (name, address, working_hours_start, working_hours_end)
VALUES ('1', '1', '8:00', '20:00');

INSERT INTO item_stores (in_stock_number, store_id, item_id)
VALUES (2, 1, 1),
       (0, 1, 2),
       (1, 1, 3);

INSERT INTO orders (id, status, total_price, deposit, rental_period_start, rental_period_end, address, latitude,
                    longitude, company_name, user_id)
VALUES (56, 'accepted', 5343, 0, '2004-10-19 20:00:00+03', '2004-10-30 18:00:00+03', 'address', '55.04868',
        '82.988786', 'vd', 103),
       (57, 'accepted', 5343, 0, '2004-10-19 20:00:00+03', '2004-10-30 18:00:00+03', 'address', '54.98254',
        '82.814378', 'vd',  103),
       (58, 'accepted', 5343, 0, '2004-10-19 20:00:00+03', '2004-10-30 18:00:00+03', 'address', '54.96244',
        '82.885103', 'vd',  103),
       (59, 'accepted', 5343, 0, '2004-10-19 20:00:00+03', '2004-10-30 18:00:00+03', 'address', '54.988017',
        '83.015966', 'vd',  103),
       (60, 'accepted', 5343, 0, '2004-10-19 20:00:00+03', '2004-10-30 18:00:00+03', 'address', '54.849023',
        '83.109914', 'vd',  103),
       (61, 'accepted', 5343, 0, '2004-10-19 20:00:00+03', '2004-10-30 18:00:00+03', 'address', '54.864174',
        '83.092518', 'vd',  103),
       (62, 'accepted', 5343, 0, '2004-10-19 20:00:00+03', '2004-10-30 18:00:00+03', 'address', '54.850213',
        '83.046704', 'vd',  103),
       (63, 'accepted', 5343, 0, '2004-10-19 20:00:00+03', '2004-10-30 18:00:00+03', 'address', '54.837411',
        '83.112056', 'vd',  103);

INSERT INTO deliveries (id, time_start, time_end, method, order_id, deliveryman_id)
VALUES (101, '2004-10-19 10:00:00+03', '2004-10-19 12:00:00+03', 'car', 56, null),
       (102, '2004-10-19 13:00:00+03', '2004-10-19 15:00:00+03', 'car', 57, null),
       (103, '2004-10-19 16:00:00+03', '2004-10-19 18:00:00+03', 'car', 58, null),
       (104, '2004-10-19 18:00:00+03', '2004-10-19 20:00:00+03', 'car', 59, null),
       (105, '2004-10-19 11:00:00+03', '2004-10-19 13:00:00+03', 'car', 60, null),
       (106, '2004-10-19 14:00:00+03', '2004-10-19 16:00:00+03', 'car', 61, null),
       (107, '2004-10-19 16:00:00+03', '2004-10-19 18:00:00+03', 'car', 62, null),
       (108, '2004-10-19 16:00:00+03', '2004-10-19 18:00:00+03', 'car', 63, null);

