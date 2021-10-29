drop table users;
drop table whitelisted_admins;
drop table admin_configs;
drop table allowed_apps;
drop table app_votes;
drop table registered_apps;
drop table stream_sessions;
drop table smart_otps;

CREATE TABLE IF NOT EXISTS users (
    id INTEGER,
    wallet_address TEXT PRIMARY KEY,
    nonce TEXT,
    status INTEGER,
    max_connections INTEGER,
    cur_unreleased_balance INTEGER,
    machine TEXT,
    location TEXT,
    name TEXT,
    hourly_rate INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS whitelisted_admins (
    id INTEGER,
    wallet_address TEXT PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS admin_configs (
    id INTEGER,
    hourly_rate INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS allowed_apps (
    id INTEGER,
    app_name TEXT PRIMARY KEY,
    publisher TEXT,
    image_url TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS app_votes (
    id INTEGER,
    app_name TEXT,
    wallet_address TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    PRIMARY KEY (app_name, wallet_address)
);

CREATE TABLE IF NOT EXISTS registered_apps (
    id INTEGER,
    wallet_address TEXT,
    app_path TEXT,
    app_name TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    PRIMARY KEY (wallet_address, app_name)
);

CREATE TABLE IF NOT EXISTS stream_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    stream_status INTEGER,
    max_duration INTEGER,
    total_duration INTEGER,
    accum_charge INTEGER,
    client_wallet_address TEXT,
    host_wallet_address TEXT,
    app_name TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS smart_otps (
    id INTEGER,
    wallet_address TEXT,
    otp TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

INSERT INTO whitelisted_admins (id, wallet_address) VALUES (1, '0x3bfe0fcaecdb3ad87c786447f10d57bd0c6cb842');
