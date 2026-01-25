-- Create tags table
CREATE TABLE IF NOT EXISTS tags
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ  NULL
);

-- Create index for soft delete queries
CREATE INDEX idx_tags_deleted_at ON tags (deleted_at);

-- Create index for name lookups
CREATE INDEX idx_tags_name ON tags (name);
