-- +migrate Up
CREATE TABLE accounts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY NOT NULL,
    fullname VARCHAR(255) NULL,
    sort_bio VARCHAR(255) NULL,
    gender ENUM('male','female','others') DEFAULT 'others',
    picture_url VARCHAR(255) DEFAULT NULL,
    username VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    role ENUM('admin','member') DEFAULT 'member',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS accounts;
