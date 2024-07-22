CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: users
CREATE TABLE users (
                       user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       name VARCHAR(255) NOT NULL,
                       gender VARCHAR(6) CHECK (gender IN ('male', 'female')) NOT NULL,
                       telp VARCHAR(15),
                       birthdate DATE NOT NULL ,
                       address TEXT,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: photos
CREATE TABLE photos (
                        id uuid PRIMARY KEY default uuid_generate_v4(),
                        url VARCHAR(255) NOT NULL,
                        owner_id uuid,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: roles
CREATE TABLE roles (
                       id serial PRIMARY KEY ,
                       user_id uuid REFERENCES users(user_id),
                       role VARCHAR(50) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP DEFAULT NULL
);


-- Table: devices
CREATE TABLE devices (
                         id SERIAL PRIMARY KEY,
                         user_id uuid REFERENCES users(user_id),
                         device_id VARCHAR(255) NOT NULL,
                         last_login TIMESTAMP,
                         user_agent TEXT,
                         platform VARCHAR(50),
                         vendor VARCHAR(50),
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: otps
CREATE TABLE otp(
                    id SERIAL PRIMARY KEY,
                    user_id uuid REFERENCES users(user_id),
                    otp VARCHAR(10) NOT NULL,
                    expired_date TIMESTAMP NOT NULL,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: oauths
CREATE TABLE oauths (
                        id SERIAL PRIMARY KEY,
                        email VARCHAR(255) UNIQUE NOT NULL,
                        password VARCHAR(255) NOT NULL,
                        is_enabled BOOLEAN DEFAULT FALSE,
                        username VARCHAR(255) UNIQUE,
                        user_id uuid REFERENCES users(user_id),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: suppliers
CREATE TABLE suppliers (
                           id uuid PRIMARY KEY default uuid_generate_v4(),
                           name VARCHAR(255) NOT NULL,
                           contact_info VARCHAR(255),
                           address TEXT,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: categories
CREATE TABLE categories (
                            id uuid PRIMARY KEY default uuid_generate_v4(),
                            name VARCHAR(255) NOT NULL,
                            description TEXT,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: brands
CREATE TABLE brands (
                        id SERIAL PRIMARY KEY ,
                        name VARCHAR(255) NOT NULL,
                        description TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: products
CREATE TABLE products (
                          id uuid PRIMARY KEY default uuid_generate_v4(),
                          product_name VARCHAR(255) NOT NULL,
                          sell_price NUMERIC NOT NULL,
                          call_name TEXT,
                          admin_id uuid REFERENCES users(user_id),
                          category_id uuid REFERENCES categories(id),
                          brand_id INTEGER REFERENCES brands(id),
                          supplier_id uuid REFERENCES suppliers(id),
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: orders
CREATE TABLE orders (
                        id uuid PRIMARY KEY default uuid_generate_v4(),
                        admin_id uuid REFERENCES users(user_id),
                        buyer_id uuid REFERENCES users(user_id),
                        total NUMERIC NOT NULL,
                        tax NUMERIC,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: order_items
CREATE TABLE order_items (
                             id SERIAL PRIMARY KEY,
                             order_id uuid REFERENCES orders(id),
                             product_id uuid REFERENCES products(id),
                             price NUMERIC NOT NULL,
                             quantity INTEGER NOT NULL,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: inventories
CREATE TABLE inventories (
                             id SERIAL PRIMARY KEY,
                             product_id uuid REFERENCES products(id),
                             total_quantity INTEGER NOT NULL,
                             last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             deleted_at TIMESTAMP DEFAULT NULL
);

-- Table: inventory_transactions
CREATE TABLE inventory_transactions (
                                        id SERIAL PRIMARY KEY,
                                        transaction_type VARCHAR(4) CHECK (transaction_type IN ('IN', 'OUT')) NOT NULL,
                                        price_modal NUMERIC,
                                        total_quantity INTEGER NOT NULL,
                                        last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        admin_id uuid REFERENCES users(user_id),
                                        supplier_id uuid REFERENCES suppliers(id),
                                        keterangan TEXT,
                                        product_id uuid REFERENCES products(id),
                                        deleted_at TIMESTAMP DEFAULT NULL
);

-- Indexing for performance improvement
CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_products_name ON products(product_name);
CREATE INDEX idx_orders_total ON orders(total);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);

CREATE OR REPLACE FUNCTION update_order_total()
    RETURNS TRIGGER AS $$
BEGIN
    UPDATE orders
    SET total = (
        SELECT SUM(price * quantity)
        FROM order_items
        WHERE order_id = NEW.order_id
    )
    WHERE id = NEW.order_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_order_total
    AFTER INSERT OR UPDATE OR DELETE
    ON order_items
    FOR EACH ROW
EXECUTE FUNCTION update_order_total();


