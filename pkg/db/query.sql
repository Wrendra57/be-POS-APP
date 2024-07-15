-- Table: users
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    gender VARCHAR(6) CHECK (gender IN ('male', 'female')) NOT NULL,
    telp VARCHAR(15) UNIQUE,
    birthdate DATE,
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: photos
CREATE TABLE photos (
    id SERIAL PRIMARY KEY,
    url VARCHAR(255) NOT NULL,
    user_id INTEGER REFERENCES users(user_id)
);

-- Table: roles
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    role VARCHAR(50) NOT NULL
);

-- Table: devices
CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    device_id VARCHAR(255) NOT NULL,
    last_login TIMESTAMP,
    user_agent TEXT,
    platform VARCHAR(50),
    vendor VARCHAR(50)
);

-- Table: otps
CREATE TABLE otps (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    otp VARCHAR(10) NOT NULL,
    expired_date TIMESTAMP NOT NULL
);

-- Table: oauths
CREATE TABLE oauths (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN DEFAULT TRUE,
    username VARCHAR(255) UNIQUE,
    user_id INTEGER REFERENCES users(user_id)
);

-- Table: suppliers
CREATE TABLE suppliers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    contact_info VARCHAR(255),
    address TEXT
);

-- Table: categories
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Table: brands
CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Table: products
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    sell_price NUMERIC NOT NULL,
    call_name VARCHAR(255),
    admin_id INTEGER REFERENCES users(user_id),
    category_id INTEGER REFERENCES categories(id),
    brand_id INTEGER REFERENCES brands(id),
    supplier_id INTEGER REFERENCES suppliers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: orders
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    admin_id INTEGER REFERENCES users(user_id),
    buyer_id INTEGER REFERENCES users(user_id),
    total NUMERIC NOT NULL,
    tax NUMERIC,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: order_items
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    product_id INTEGER REFERENCES products(id),
    price NUMERIC NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: inventories
CREATE TABLE inventories (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id),
    total_quantity INTEGER NOT NULL,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: inventory_transactions
CREATE TABLE inventory_transactions (
    id SERIAL PRIMARY KEY,
    transaction_type VARCHAR(4) CHECK (transaction_type IN ('IN', 'OUT')) NOT NULL,
    price_modal NUMERIC,
    total_quantity INTEGER NOT NULL,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    admin_id INTEGER REFERENCES users(user_id),
    supplier_id INTEGER REFERENCES suppliers(id),
    keterangan TEXT,
    product_id INTEGER REFERENCES products(id)
);

-- Indexing for performance improvement
CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_products_name ON products(product_name);
CREATE INDEX idx_orders_total ON orders(total);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);