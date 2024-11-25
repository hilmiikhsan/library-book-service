package models

import (
	"time"

	"github.com/google/uuid"
)

type BookBorrowed struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	BookID       uuid.UUID `db:"book_io"`
	BorrowedDate time.Time `db:"borrowed_date"`
	DueDate      time.Time `db:"due_date"`
	ReturnedDate time.Time `db:"returned_date"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
