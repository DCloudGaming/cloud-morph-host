drop table users;
drop table host_configs;
drop table registered_apps;
drop table stream_sessions;

-- Please note INTEGER is LONG too in SQLite

CREATE TABLE IF NOT EXISTS users (
    id INTEGER,
    wallet_address TEXT PRIMARY KEY,
    nonce TEXT,
    status INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS host_configs (
    id INTEGER,
    wallet_address TEXT PRIMARY KEY,
    max_connections INTEGER,
    cur_unreleased_balance INTEGER,
    hourly_rate INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
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

-- TODO: Change id to uuid
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