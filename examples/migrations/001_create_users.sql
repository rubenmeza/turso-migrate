-- Migration: Create users table
-- Created: 2024-01-05 10:00:00

-- ==== UP ====
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- ==== DOWN ====
DROP INDEX idx_users_email;
DROP TABLE users;