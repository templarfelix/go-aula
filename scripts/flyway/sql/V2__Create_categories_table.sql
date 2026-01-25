-- Create categories table
CREATE TABLE IF NOT EXISTS categories
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ  NULL
);

-- Create index for soft delete queries
CREATE INDEX idx_categories_deleted_at ON categories (deleted_at);

-- Create index for name lookups
CREATE INDEX idx_categories_name ON categories (name);
