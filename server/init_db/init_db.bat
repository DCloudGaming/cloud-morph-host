drop table user;

CREATE TABLE IF NOT EXISTS user (
    wallet_address TEXT PRIMARY KEY,
    nonce TEXT
);
