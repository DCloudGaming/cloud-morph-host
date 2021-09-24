drop table users;

CREATE TABLE IF NOT EXISTS users (
    id INTEGER,
    wallet_address TEXT PRIMARY KEY,
    nonce TEXT,
    status INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
