-- Пользователи системы
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'manager', 'waiter', 'cook')),
    photoKey VARCHAR(20) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Столы
CREATE TABLE tables (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('busy', 'reserve', 'free')) DEFAULT 'free'
);
-- Категории блюд
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
    dish_id INT NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    ingredient_id INT NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    qty_per_dish NUMERIC(10, 2) NOT NULL CHECK (qty_per_dish > 0),
    PRIMARY KEY (dish_id, ingredient_id)
);
-- Заказы
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    waiter_id INT NOT NULL REFERENCES users(id),
    table_number INT NOT NULL REFERENCES tables(id),
    status VARCHAR(20) NOT NULL CHECK (
        status IN ('new', 'in_progress', 'ready', 'paid')
    ) DEFAULT 'new',
    total NUMERIC(10, 2) DEFAULT 0 CHECK (total >= 0),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Позиции в заказе
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    dish_id INT NOT NULL REFERENCES dishes(id),
    qty INT NOT NULL CHECK (qty > 0),
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    notes TEXT
);
-- Поставки ингредиентов
CREATE TABLE supplies (
    id SERIAL PRIMARY KEY,
    ingredient_id INT NOT NULL REFERENCES ingredients(id),
    qty NUMERIC(10, 2) NOT NULL CHECK (qty > 0),
    supplier_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Индексы для производительности
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_waiter_id ON orders(waiter_id);
CREATE INDEX idx_orders_table_number ON orders(table_number);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_dish_id ON order_items(dish_id);
CREATE INDEX idx_dishes_category_id ON dishes(category_id);
CREATE INDEX idx_dishes_is_active ON dishes(is_active);
CREATE INDEX idx_supplies_ingredient_id ON supplies(ingredient_id);
CREATE INDEX idx_supplies_created_at ON supplies(created_at);