-- Пользователи системы
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (
        role IN (
            'admin',
            'manager',
            'waiter',
            'cook'
        )
    ),
    photokey VARCHAR(20) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Столы
CREATE TABLE tables (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL CHECK (
        status IN ('busy', 'reserve', 'free')
    ) DEFAULT 'free'
);
-- Категории блюд
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);
-- Блюда
CREATE TABLE dishes (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    photo_url TEXT,
    is_active BOOLEAN DEFAULT true
);
-- Ингредиенты на складе
CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    unit VARCHAR(20) NOT NULL,
    -- кг, литр, шт
    qty NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (qty >= 0),
    min_qty NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (min_qty >= 0)
);
-- Связь блюд и ингредиентов
CREATE TABLE dish_ingredients (
    dish_id INT NOT NULL REFERENCES dishes (id) ON DELETE CASCADE,
    ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    qty_per_dish NUMERIC(10, 2) NOT NULL CHECK (qty_per_dish > 0),
    PRIMARY KEY (dish_id, ingredient_id)
);
-- Заказы
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    waiter_id INT NOT NULL REFERENCES users (id),
    table_number INT NOT NULL REFERENCES tables (id),
    status VARCHAR(20) NOT NULL CHECK (
        status IN (
            'new',
            'in_progress',
            'ready',
            'paid'
        )
    ) DEFAULT 'new',
    total NUMERIC(10, 2) DEFAULT 0 CHECK (total >= 0),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Позиции в заказе
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    dish_id INT NOT NULL REFERENCES dishes (id),
    qty INT NOT NULL CHECK (qty > 0),
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    notes TEXT
);
-- Поставки ингредиентов
CREATE TABLE supplies (
    id SERIAL PRIMARY KEY,
    ingredient_id INT NOT NULL REFERENCES ingredients (id),
    qty NUMERIC(10, 2) NOT NULL CHECK (qty > 0),
    supplier_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Индексы для производительности
CREATE INDEX idx_orders_status ON orders (status);

CREATE INDEX idx_orders_waiter_id ON orders (waiter_id);

CREATE INDEX idx_orders_table_number ON orders (table_number);

CREATE INDEX idx_orders_created_at ON orders (created_at);

CREATE INDEX idx_order_items_order_id ON order_items (order_id);

CREATE INDEX idx_order_items_dish_id ON order_items (dish_id);

CREATE INDEX idx_dishes_category_id ON dishes (category_id);

CREATE INDEX idx_dishes_is_active ON dishes (is_active);

CREATE INDEX idx_supplies_ingredient_id ON supplies (ingredient_id);

CREATE INDEX idx_supplies_created_at ON supplies (created_at);

-- === USERS TABLE SEED DATA ===
INSERT INTO
    users (
        username,
        password_hash,
        role,
        photokey,
        is_active
    )
VALUES
    -- Пароль: admin123
    (
        'admin',
        '$2a$12$ZZwelVZlMwCeGd0sX019DODX0tFvFL0jzNStY6eTGgBtJQ3IsQuQq',
        'admin',
        'photo_admin',
        true
    ),

-- Пароль: manager123
(
    'manager',
    '$2a$12$r75Az7YZ3snhTXb6f8T/WudNNx/jROr/cTLswp.rYVobSs58ZdIF',
    'manager',
    'photo_manager',
    true
),

-- Пароль: waiter123
(
    'waiter',
    '$2a$10$Wgtn/qv.qUmcD8e8yEKzleJ6lwJzn/9E1pmCUYH/tHdV2ujA4R1Pa',
    'waiter',
    'photo_waiter',
    true
),

-- Пароль: cook123
(
    'cook',
    '$2a$10$Jji0mDGuA0F1L2iMnC9ybeE.4hjZ0h9u8eHHeIKrDml.nmZ0VUZLy',
    'cook',
    'photo_cook',
    true
);

-- === EXISTING USERS TABLE (unchanged) ===
-- см. твою вставку выше

-- === TABLES SEED DATA ===
INSERT INTO
    tables (name, status)
VALUES ('Table 1', 'free'),
    ('Table 2', 'busy'),
    ('Table 3', 'reserve'),
    ('Table 4', 'free'),
    ('Table 5', 'free');

-- === CATEGORIES SEED DATA ===
INSERT INTO
    categories (name)
VALUES ('Salads'),
    ('Soups'),
    ('Main Dishes'),
    ('Desserts'),
    ('Drinks');

-- === DISHES SEED DATA ===
INSERT INTO
    dishes (
        category_id,
        name,
        description,
        price,
        photo_url,
        is_active
    )
VALUES (
        1,
        'Caesar Salad',
        'Classic Caesar with chicken and parmesan',
        2500,
        'dishes/caesar.jpg',
        true
    ),
    (
        1,
        'Greek Salad',
        'Fresh vegetables, feta cheese and olives',
        2200,
        'dishes/greek.jpg',
        true
    ),
    (
        2,
        'Tomato Soup',
        'Homemade tomato soup with basil',
        1800,
        'dishes/tomato_soup.jpg',
        true
    ),
    (
        3,
        'Beef Steak',
        'Grilled beef steak with sauce',
        7500,
        'dishes/steak.jpg',
        true
    ),
    (
        3,
        'Chicken Curry',
        'Spicy chicken curry with rice',
        6500,
        'dishes/curry.jpg',
        true
    ),
    (
        4,
        'Cheesecake',
        'Classic New York cheesecake',
        3000,
        'dishes/cheesecake.jpg',
        true
    ),
    (
        5,
        'Cappuccino',
        'Espresso with milk foam',
        1500,
        'dishes/cappuccino.jpg',
        true
    ),
    (
        5,
        'Orange Juice',
        'Freshly squeezed orange juice',
        1200,
        'dishes/orange_juice.jpg',
        true
    );

-- === INGREDIENTS SEED DATA ===
INSERT INTO
    ingredients (name, unit, qty, min_qty)
VALUES ('Chicken Breast', 'kg', 10, 2),
    ('Beef', 'kg', 8, 2),
    ('Rice', 'kg', 15, 5),
    ('Lettuce', 'kg', 5, 1),
    ('Tomato', 'kg', 7, 2),
    ('Cucumber', 'kg', 6, 2),
    ('Cheese', 'kg', 4, 1),
    ('Flour', 'kg', 12, 3),
    ('Sugar', 'kg', 10, 3),
    ('Milk', 'liter', 20, 5),
    ('Coffee Beans', 'kg', 6, 2),
    ('Orange', 'kg', 8, 3);

------------------------------------------------------------
-- === DISH INGREDIENT FORMULAS (REALISTIC RECIPES) ===
-- Полное заполнение dish_ingredients по именам блюд и ингридиентов
------------------------------------------------------------

-- Очистка старых связей (если скрипт запускается повторно)
DELETE FROM dish_ingredients;

-- Вставляем формулы
INSERT INTO dish_ingredients (dish_id, ingredient_id, qty_per_dish)
SELECT d.id, i.id, v.qty
FROM (
    -- Caesar Salad
    VALUES
        ('Caesar Salad',  'Lettuce',        0.15),
        ('Caesar Salad',  'Chicken Breast', 0.20),
        ('Caesar Salad',  'Cheese',         0.05),

-- Greek Salad
('Greek Salad', 'Tomato', 0.20),
(
    'Greek Salad',
    'Cucumber',
    0.15
),
('Greek Salad', 'Cheese', 0.05),
(
    'Greek Salad',
    'Lettuce',
    0.10
),

-- Tomato Soup
('Tomato Soup', 'Tomato', 0.30), ('Tomato Soup', 'Rice', 0.05),

-- Beef Steak
('Beef Steak', 'Beef', 0.35),

-- Chicken Curry
(
    'Chicken Curry',
    'Chicken Breast',
    0.25
),
(
    'Chicken Curry',
    'Tomato',
    0.10
),
('Chicken Curry', 'Rice', 0.20),
('Chicken Curry', 'Milk', 0.10),

-- Cheesecake
('Cheesecake', 'Flour', 0.10),
('Cheesecake', 'Sugar', 0.05),
('Cheesecake', 'Milk', 0.15),
('Cheesecake', 'Cheese', 0.20),

-- Cappuccino
(
    'Cappuccino',
    'Coffee Beans',
    0.02
),
('Cappuccino', 'Milk', 0.20),

-- Orange Juice


('Orange Juice',  'Orange',         0.25)

) AS v(dish_name, ingredient_name, qty)
JOIN dishes d ON d.name = v.dish_name
JOIN ingredients i ON i.name = v.ingredient_name;

-- === DISH_INGREDIENTS SEED DATA ===
INSERT INTO
    dish_ingredients (
        dish_id,
        ingredient_id,
        qty_per_dish
    )
VALUES (1, 4, 0.2), -- Caesar Salad - Lettuce
    (1, 1, 0.15), -- Chicken Breast
    (1, 7, 0.05), -- Cheese
    (2, 5, 0.1), -- Greek Salad - Tomato
    (2, 6, 0.1), -- Cucumber
    (2, 7, 0.05), -- Cheese
    (3, 5, 0.15), -- Tomato Soup - Tomato
    (4, 2, 0.3), -- Beef Steak - Beef
    (5, 1, 0.25), -- Chicken Curry - Chicken
    (5, 3, 0.1), -- Rice
    (6, 8, 0.1), -- Cheesecake - Flour
    (6, 9, 0.05), -- Sugar
    (6, 10, 0.1), -- Milk
    (7, 10, 0.2), -- Cappuccino - Milk
    (7, 11, 0.05), -- Coffee Beans
    (8, 12, 0.3);
-- Orange Juice - Orange

-- === ORDERS SEED DATA ===
INSERT INTO
    orders (
        waiter_id,
        table_number,
        status,
        total,
        notes
    )
VALUES (
        3,
        2,
        'in_progress',
        9000,
        'Customer allergic to nuts'
    ),
    (3, 3, 'new', 6500, NULL),
    (
        3,
        1,
        'paid',
        12000,
        'VIP guest'
    );

-- === ORDER_ITEMS SEED DATA ===
INSERT INTO
    order_items (
        order_id,
        dish_id,
        qty,
        price,
        notes
    )
VALUES (1, 1, 1, 2500, NULL),
    (1, 4, 1, 7500, 'Medium rare'),
    (2, 5, 1, 6500, 'Less spicy'),
    (3, 4, 1, 7500, NULL),
    (
        3,
        6,
        1,
        3000,
        'Extra topping'
    );

-- === SUPPLIES SEED DATA ===
INSERT INTO
    supplies (
        ingredient_id,
        qty,
        supplier_name
    )
VALUES (1, 5, 'FreshMeat Co'),
    (4, 3, 'GreenFarm'),
    (5, 4, 'VeggieWorld'),
    (10, 10, 'DairyBest'),
    (11, 2, 'CoffeePlanet'),
    (12, 5, 'CitrusHouse');