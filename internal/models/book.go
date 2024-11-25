package models

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID            uuid.UUID `db:"id"`
	Title         string    `db:"title"`
	AuthorID      uuid.UUID `db:"author_id"`
	CategoryID    uuid.UUID `db:"category_id"`
	Description   string    `db:"description"`
	Isbn          *string   `db:"isbn"`
	PublishedDate time.Time `db:"published_date"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
