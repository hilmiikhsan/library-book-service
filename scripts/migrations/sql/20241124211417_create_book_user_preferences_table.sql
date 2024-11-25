-- +goose Up
-- +goose StatementBegin
CREATE TABLE book_user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    preferred_category UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster query by user_id
CREATE INDEX idx_book_user_preferences_user_id 
ON book_user_preferences(user_id);

-- Index for faster query by preferred_category
CREATE INDEX idx_book_user_preferences_category_id 
ON book_user_preferences(preferred_category);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_book_user_preferences_user_id;
DROP INDEX IF EXISTS idx_book_user_preferences_category_id;
DROP TABLE IF EXISTS book_user_preferences;
-- +goose StatementEnd
