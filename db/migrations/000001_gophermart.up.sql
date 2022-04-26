CREATE TABLE IF NOT EXISTS users
(
    login character varying PRIMARY KEY,
    password character varying NOT NULL
);
CREATE TABLE IF NOT EXISTS orders
(
    "number" character varying PRIMARY KEY,
    "user" character varying NOT NULL,
    status character varying NOT NULL,
    accrual numeric NOT NULL DEFAULT 0,
    uploaded_at timestamp with time zone default CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS balances
(
    processed_at timestamp without time zone,
    "user" character varying NOT NULL,
    "order" character varying NOT NULL,
    sum numeric NOT NULL DEFAULT 0
);
