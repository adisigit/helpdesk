DROP INDEX IF EXISTS idx_comments_author_id;

ALTER TABLE comments
    DROP COLUMN IF EXISTS author_id;
