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

INSERT INTO subscriptions (user_id)
VALUES (100),
       (101),
       (102),
       (103);

INSERT INTO deliverymen (car_capacity, working_hours_start, working_hours_end, car_id, user_id)
VALUES ('big', '10:00', '16:00', 'A123BC', 101),
       ('medium', '12:00', '20:00', 'X777XX', 102);

INSERT INTO categories (id, name, parent_id)
VALUES (1, 'Ручной инструмент', 0),
       (2, 'Электроинструменты', 0),
       (3, 'Отбойные', 1),
       (4, 'Измерительные', 1),
       (5, 'Режущие', 2),
       (6, 'Сверлющие', 2);

INSERT INTO items (name, description, price, price_deposit, is_in_stock, category_id)
VALUES ('Молоток', 'Молоток слесарный с деревянной рукояткой предназначен для ремонтных и строительных работ', 2000, 20000, true, 3),
       ('Пила', 'Точная, быстрая резка, Индукционно закаленные зубья, Прочные и жесткие зубья с трехгранной заточкой', 3000, 30000, true, 5),
       ('Кувалда', 'Кованный закаленный боек из высокоуглеродистой стали', 1500, 10000, true, 3),
       ('Топор', 'Топор с деревянной рукояткой предназначен для рубки, колки и тески древесины', 1000, 5000, true, 5),
       ('Выколотка', 'Выколотки слесарные предназначены для монтажа и демонтажа штифтов, шплинтов и т.п', 300, 3000, false, 4),
       ('Ножовка по дереву', 'Ножовка предназначена для резки древесины, древесных материалов и мягкого пластика', 2400, 12000, true, 5),
       ('Болгарка', 'Хорошо пилит, мощно', 6000, 12000, true, 5),
       ('Труборез', 'Хорошо режет, мощно', 2000, 4000, true, 5),
       ('Перфоратор', 'Хорошо сверлит, мощно', 1000, 1500, true, 4),
       ('Кернер', 'Хорошо рамечевает, мощно', 2000, 20000, true, 4),
       ('Струбцина', 'Хорошо держит, крепко', 3000, 30000, true, 4),
       ('Паяльник', 'Хорошо паяет, горячий', 4000, 12000, true, 2),
       ('Мультиметр', 'Хорошо мерит, мульти', 5000, 18000, true, 4),
       ('Зубило', 'Хорошо зубилит, зубасто', 6000, 12000, true, 5),
       ('Рубанок', 'Хорошо рубит, он рубанок', 800, 4000, true, 5),
       ('Дрель', 'Хорошо сверлит, глубоко', 3500, 30000, true, 6),
       ('Лом', 'Хорошо ломает, но не строит', 2000, 30000, true, 3),
       ('Стамеска', 'Хорошо обрабатывает, мощно', 3100, 30000, true, 5),
       ('Плоскогубцы', 'Хорошо сжимает, крепкая хватка', 4000, 40000, true, 1);

INSERT INTO infos (name, description, item_id)
VALUES ('Габариты', '120x120x120', 1),
       ('Габариты', '120x120x120', 2),
       ('Габариты', '120x120x120', 3),
       ('Габариты', '120x120x120', 4),
       ('Габариты', '120x120x120', 5),
       ('Габариты', '120x120x120', 6),
       ('Габариты', '120x120x120', 7),
       ('Габариты', '120x120x120', 8),
       ('Вес', '1 кг', 6),
       ('Вес', '2 кг', 7),
       ('Вес', '3 кг', 9),
       ('Вес', '1 кг', 10),
       ('Вес', '4 кг', 11),
       ('Вес', '5 кг', 13),
       ('Вес', '6 кг', 14),
       ('Вес', '7 кг', 12),
       ('Вес', '12 кг', 15),
       ('Вес', '8 кг', 18),
       ('Материал', 'Дерево', 1),
       ('Материал', 'Дерево', 17),
       ('Материал', 'Метал', 18),
       ('Материал', 'Пластик', 19),
       ('Материал', 'Дерево', 16);

INSERT INTO item_images (image, item_id)
VALUES ('hammer.jpg', 1),
       ('saw.jpg', 2),
       ('kuvalda.jpg', 3),
       ('axe.jpg', 4),
       ('vikolotka.jpg', 5),
       ('nozhovka.jpg', 6),
       ('bolgarka.jpeg', 7),
       ('truborez.jpg', 8),
       ('perf.jpg', 9),
       ('kerner.jpg', 10),
       ('strubcina.jpg', 11),
       ('payalnik.jpg', 12),
       ('multimetr.jpg', 13),
       ('zubilo.jpg', 14),
       ('rubanok.jpg', 15),
       ('drell.jpeg', 16),
       ('lom.jpg', 17),
       ('stameska.jpg', 18),
       ('set.png', 19);

INSERT INTO stores (name, address, latitude, longitude, working_hours_start, working_hours_end)
VALUES ('1', 'Россия Новосибирск Пирогова 1', '54.843072', '83.090792', '8:00', '20:00');

INSERT INTO item_stores (in_stock_number, store_id, item_id)
VALUES (58, 1, 1),
       (10, 1, 2),
       (10, 1, 3),
       (10, 1, 4),
       (0, 1, 5),
       (10, 1, 6),
       (10, 1, 7),
       (10, 1, 8),
       (10, 1, 9),
       (10, 1, 10),
       (10, 1, 11),
       (10, 1, 12),
       (10, 1, 13),
       (10, 1, 14),
       (10, 1, 15),
       (10, 1, 16),
       (1, 1, 17),
       (10, 1, 18),
       (0, 1, 19);

