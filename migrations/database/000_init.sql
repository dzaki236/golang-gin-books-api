-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    modified_at TIMESTAMP,
    modified_by VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    modified_at TIMESTAMP,
    modified_by VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    image_url VARCHAR(2048) NOT NULL DEFAULT '',
    release_year INTEGER NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    total_page INTEGER NOT NULL CHECK (total_page > 0),
    thickness VARCHAR(20) NOT NULL,
    category_id INTEGER NOT NULL REFERENCES categories (id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    modified_at TIMESTAMP,
    modified_by VARCHAR(100)
);
-- +migrate StatementEnd
-- +migrate Down

DROP TABLE IF EXISTS books;

DROP TABLE IF EXISTS categories;

DROP TABLE IF EXISTS users;
