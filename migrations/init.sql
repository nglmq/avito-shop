CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS balances (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) REFERENCES users(username),
    balance INT NOT NULL DEFAULT 1000 CHECK (balance >= 0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) REFERENCES users(username),
    item_name VARCHAR(255) NOT NULL,
    amount INT NOT NULL DEFAULT 1,
    total_price INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sender_username VARCHAR(255) REFERENCES users(username),
    receiver_username VARCHAR(255) REFERENCES users(username),
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_transactions_sender_username ON transactions(sender_username);
CREATE INDEX IF NOT EXISTS idx_transactions_receiver_username ON transactions(receiver_username);