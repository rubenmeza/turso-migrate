-- Migration: Add posts table with foreign key to users
-- Created: 2024-01-05 10:15:00

-- ==== UP ====
CREATE TABLE posts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    published BOOLEAN NOT NULL DEFAULT false,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_published ON posts(published);

-- ==== DOWN ====
DROP INDEX idx_posts_published;
DROP INDEX idx_posts_user_id;
DROP TABLE posts;