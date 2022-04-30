DROP TABLE IF EXISTS assets CASCADE;

CREATE TABLE assets (
  id SERIAL PRIMARY KEY NOT NULL,
  token_name VARCHAR(255) NOT NULL,
  quantity DOUBLE PRECISION NOT NULL,
  wallet_id INTEGER REFERENCES wallets(id) ON DELETE CASCADE
);