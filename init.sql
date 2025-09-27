-- Пользователи системы
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) CHECK (role IN ('admin', 'manager', 'waiter', 'cook')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Категории блюд (салаты, супы и т.д.)
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);
-- Блюда
CREATE TABLE dishes (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    photo_url TEXT
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
-- Связь блюд и ингредиентов (многие-ко-многим)
CREATE TABLE dish_ingredients (
    dish_id INT NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    ingredient_id INT NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    qty_per_dish NUMERIC(10, 2) NOT NULL CHECK (qty_per_dish > 0),
    PRIMARY KEY (dish_id, ingredient_id)
);
-- Заказы
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    waiter_id INT NOT NULL REFERENCES users(id),
    table_number INT NOT NULL CHECK (table_number > 0),
    status VARCHAR(20) NOT NULL CHECK (
        status IN ('new', 'in_progress', 'ready', 'paid')
    ),
    total NUMERIC(10, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Позиции в заказе
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    dish_id INT NOT NULL REFERENCES dishes(id),
    qty INT NOT NULL CHECK (qty > 0),
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0)
);
-- Поставки ингредиентов
CREATE TABLE supplies (
    id SERIAL PRIMARY KEY,
    ingredient_id INT NOT NULL REFERENCES ingredients(id),
    qty NUMERIC(10, 2) NOT NULL CHECK (qty > 0),
    supplier_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);