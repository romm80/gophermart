CREATE TABLE IF NOT EXISTS users
(
    id serial PRIMARY KEY,
    login character varying NOT NULL UNIQUE,
    password character varying NOT NULL
);
CREATE TABLE IF NOT EXISTS orders
(
    id serial PRIMARY KEY,
    "number" character varying NOT NULL UNIQUE,
    user_id serial REFERENCES users (id),
    status character varying NOT NULL,
    accrual numeric NOT NULL DEFAULT 0,
    uploaded_at timestamp with time zone default CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS balances
(
    id serial PRIMARY KEY,
    processed_at timestamp without time zone,
    user_id serial REFERENCES users (id),
    order_id serial REFERENCES orders (id),
    sum numeric NOT NULL DEFAULT 0
);
