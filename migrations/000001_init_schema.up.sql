CREATE TABLE users (
  id INTEGER NOT NULL,
  created_at TIMESTAMP,
  name TEXT NOT NULL,
  email TEXT NOT NULL
);

CREATE TABLE products (
  id INTEGER NOT NULL,
  price NUMERIC NOT NULL,
  stock INTEGER NOT NULL,
  created_at TIMESTAMP,
  name TEXT NOT NULL
);

CREATE TABLE orders (
  id INTEGER NOT NULL,
  total_price NUMERIC NOT NULL,
  created_at TIMESTAMP,
  user_id INTEGER NOT NULL,
  user_email TEXT
);

CREATE TABLE history (
  id INTEGER NOT NULL,
  order_id INTEGER,
  product_id INTEGER,
  quantity INTEGER,
  price NUMERIC NOT NULL
);

CREATE TABLE transactions (
  id INTEGER NOT NULL,
  order_id INTEGER,
  amount NUMERIC NOT NULL,
  created_at TIMESTAMP,
  status TEXT NOT NULL
);

CREATE TABLE order_items (
  id INTEGER NOT NULL,
  order_id INTEGER NOT NULL,
  product_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL,
  price NUMERIC NOT NULL
);

CREATE TABLE popular_products (
  id INTEGER NOT NULL,
  product_id INTEGER,
  total_sold INTEGER NOT NULL,
  product_name VARCHAR(255)
);