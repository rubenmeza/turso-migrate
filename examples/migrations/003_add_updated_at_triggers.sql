-- Migration: Add updated_at trigger for users
-- Created: 2024-01-05 10:30:00

-- ==== UP ====
-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER users_updated_at
    AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Create trigger for posts table as well
CREATE TRIGGER posts_updated_at
    AFTER UPDATE ON posts
BEGIN
    UPDATE posts SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- ==== DOWN ====
DROP TRIGGER posts_updated_at;
DROP TRIGGER users_updated_at;