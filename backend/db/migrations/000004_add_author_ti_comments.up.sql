ALTER TABLE comments
    ADD COLUMN IF NOT EXISTS author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_comments_author_id ON comments(author_id);
