package models

import (
	"time"

	"github.com/google/uuid"
)

type BookStock struct {
	ID             uuid.UUID `db:"id"`
	BookID         uuid.UUID `db:"book_id"`
	BookTitle      string    `db:"book_title"`
	TotalStock     int       `db:"total_stock"`
	AvailableStock int       `db:"available_stock"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
