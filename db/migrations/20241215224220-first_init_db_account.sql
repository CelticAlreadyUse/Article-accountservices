-- +migrate Up
CREATE TABLE accounts (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    username VARCHAR(255) NOT NULL UNIQUE,
    gender ENUM('male','female','others') DEFAULT 'others',
    password_hash VARCHAR(255),
    picture_url VARCHAR(255) DEFAULT NULL,
    display_name VARCHAR(50),
    short_bio VARCHAR(160) NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    email_verify BOOLEAN DEFAULT FALSE,
    role ENUM('admin','member') DEFAULT 'member',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);
-- +migrate Down
DROP TABLE IF EXISTS accounts;
