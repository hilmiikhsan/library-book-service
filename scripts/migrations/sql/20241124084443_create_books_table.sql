-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    author_id UUID NOT NULL, -- Author reference, no FK enforced
    category_id UUID NOT NULL, -- Category reference, no FK enforced
    description TEXT,
    isbn VARCHAR(20) UNIQUE NOT NULL,
    published_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- GIN index for full-text search
CREATE INDEX idx_books_title_description_gin 
ON books USING gin(to_tsvector('english', title || ' ' || description));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_books_title_description_gin;
DROP TABLE IF EXISTS books;
-- +goose StatementEnd
