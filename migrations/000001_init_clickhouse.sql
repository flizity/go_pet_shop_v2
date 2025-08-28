CREATE TABLE users (
  id UInt32,
  created_at DateTime,
  name String,
  email String
) ENGINE = MergeTree() ORDER BY id;

CREATE TABLE products (
  id UInt32,
  price Decimal(10,2),
  stock UInt32,
  created_at DateTime,
  name String
) ENGINE = MergeTree() ORDER BY id;

CREATE TABLE orders (
  id UInt32,
  total_price Decimal(10,2),
  created_at DateTime,
  user_id UInt32,
  user_email String
) ENGINE = MergeTree() ORDER BY id;

CREATE TABLE history (
  id UInt32,
  order_id UInt32,
  product_id UInt32,
  quantity UInt32,
  price Decimal(10,2)
) ENGINE = MergeTree() ORDER BY id;

CREATE TABLE transactions (
  id UInt32,
  order_id UInt32,
  amount Decimal(10,2),
  created_at DateTime,
  status String
) ENGINE = MergeTree() ORDER BY id;

CREATE TABLE order_items (
  id UInt32,
  order_id UInt32,
  product_id UInt32,
  quantity UInt32,
  price Decimal(10,2)
) ENGINE = MergeTree() ORDER BY id;

CREATE TABLE popular_products (
  id UInt32,
  product_id UInt32,
  total_sold UInt32,
  product_name String
) ENGINE = MergeTree() ORDER BY id;