CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,         
    menu_item_name VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

