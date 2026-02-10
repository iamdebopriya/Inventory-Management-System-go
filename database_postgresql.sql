-- Check PostgreSQL version
SELECT version();
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE
);
-- Create products table
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_name VARCHAR(150) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    quantity INT NOT NULL CHECK (quantity >= 0),
    category_id UUID REFERENCES categories(id),
    is_active BOOLEAN DEFAULT TRUE
);
-- Create orders table
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id),
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    quantity INT NOT NULL CHECK (quantity > 0)
);
-- List all tables in the public schema
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public';

-- Insert sample data into categories and products tables
INSERT INTO categories (name) VALUES
('Electronics'),
('Books'),
('Clothing');
-- Insert sample data into products table
INSERT INTO products (product_name, price, quantity, category_id) VALUES
('Laptop', 900.00, 5, (SELECT id FROM categories WHERE name='Electronics')),
('Mouse', 25.00, 20, (SELECT id FROM categories WHERE name='Electronics')),
('Keyboard', 60.00, 10, (SELECT id FROM categories WHERE name='Electronics')),
('Headphones', 120.00, 0, (SELECT id FROM categories WHERE name='Electronics')),
('Novel', 15.00, 30, (SELECT id FROM categories WHERE name='Books')),
('Textbook', 80.00, 3, (SELECT id FROM categories WHERE name='Books')),
('T-Shirt', 20.00, 50, (SELECT id FROM categories WHERE name='Clothing')),
('Jeans', 70.00, 8, (SELECT id FROM categories WHERE name='Clothing')),
('Jacket', 150.00, 2, (SELECT id FROM categories WHERE name='Clothing')),
('Shoes', 95.00, 4, (SELECT id FROM categories WHERE name='Clothing'));
-- Simple SELECT query
SELECT * FROM products;

-- SELECT products with quantity > 0 and price > 50
SELECT *
FROM products
WHERE quantity > 0 AND price > 50;
-- JOIN products with categories to get category names
SELECT p.product_name, c.name AS category_name
FROM products p
INNER JOIN categories c ON p.category_id = c.id;
-- Aggregate query to get total inventory value per category
SELECT c.name AS category_name,
       SUM(p.price * p.quantity) AS total_inventory_value
FROM products p
JOIN categories c ON p.category_id = c.id
GROUP BY c.name;
-- Create an index on product_name column using BTREE
CREATE INDEX idx_product_name ON products USING BTREE(product_name);
-- Verify the index creation
SELECT * FROM products WHERE product_name = 'Laptop';
-- Transaction block to handle purchase
DO $$
DECLARE
    current_qty INT;
    prod_id UUID := '9e9e1892-55ce-434f-82d1-24c6d5917123';  
BEGIN
    -- Check current quantity
    SELECT quantity INTO current_qty
    FROM products
    WHERE id = prod_id;

    IF current_qty > 0 THEN
        -- Subtract 1 from stock
        UPDATE products
        SET quantity = quantity - 1
        WHERE id = prod_id;

        -- Insert into orders
        INSERT INTO orders (product_id, quantity)
        VALUES (prod_id, 1);

        RAISE NOTICE 'Purchase successful!';
    ELSE
        -- Stock is zero, rollback this block automatically
        RAISE Exception 'Out of stock! Purchase not completed.';
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        RAISE NOTICE 'Transaction failed: %', SQLERRM;
END $$;


-- Explanation about indexing BOOLEAN columns:
-- We do not create an index on columns like is_active (BOOLEAN)
-- because such columns have very low cardinality (only TRUE or FALSE).
-- An index is useful when it helps reduce the number of rows scanned.
-- With only two possible values, the database will still scan many rows,
-- so the index gives little or no performance benefit and increases
-- storage size and slows down INSERT/UPDATE operations.


-- Verify the created index on products table
SELECT indexname, indexdef
FROM pg_indexes
WHERE tablename = 'products';



--- List all databases and users
SELECT datname FROM pg_database;
SELECT usename FROM pg_user;

--- Create a new user with specific privileges
CREATE USER inventory_user WITH PASSWORD '';

GRANT CONNECT ON DATABASE inventory_db TO inventory_user;
GRANT USAGE ON SCHEMA public TO inventory_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO inventory_user;


ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO inventory_user;


GRANT USAGE ON SCHEMA public TO inventory_user;
GRANT CREATE ON SCHEMA public TO inventory_user;


ALTER TABLE public.products OWNER TO inventory_user;
ALTER TABLE public.categories OWNER TO inventory_user;
ALTER TABLE public.orders OWNER TO inventory_user;



