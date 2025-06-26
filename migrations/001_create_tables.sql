-- +migrate Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    wallet_id VARCHAR(64) UNIQUE NOT NULL,
    username VARCHAR(64) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    currency VARCHAR(8) NOT NULL,
    balance NUMERIC(20,2) DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE bets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    provider_tx_id VARCHAR(128) UNIQUE NOT NULL,
    amount NUMERIC(20,2) NOT NULL,
    status VARCHAR(16) NOT NULL,
    withdrawn_tx_id VARCHAR(128),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    bet_id INTEGER REFERENCES bets(id),
    type VARCHAR(16) NOT NULL,
    amount NUMERIC(20,2) NOT NULL,
    old_balance NUMERIC(20,2) NOT NULL,
    new_balance NUMERIC(20,2) NOT NULL,
    status VARCHAR(16) NOT NULL,
    provider_tx_id VARCHAR(128),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +migrate Down
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS bets;
DROP TABLE IF EXISTS users; 