package book_borrowed

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type BookBorrowedRepository struct {
	DB     *sqlx.DB
	Logger *logrus.Logger
}

func (r *BookBorrowedRepository) InsertNewBookBorrowed(ctx context.Context, tx *sql.Tx, bookBorrowed *models.BookBorrowed) error {
	_, err := tx.ExecContext(ctx, r.DB.Rebind(queryInsertNewBookBorrowed),
		bookBorrowed.UserID,
		bookBorrowed.BookID,
		bookBorrowed.DueDate,
	)
	if err != nil {
		r.Logger.Error("repo::InsertNewBookBorrowed - Failed to insert new book borrowed : ", err)
		return err
	}

	return nil
}

func (r *BookBorrowedRepository) ValidateBookBorrowed(ctx context.Context, tx *sql.Tx, bookID, userID string) error {
	var (
		count int
	)

	err := tx.QueryRowContext(ctx, r.DB.Rebind(queryValidateBookBorrowed), bookID, userID).Scan(&count)
	if err != nil {
		r.Logger.Error("repo::ValidateBookBorrowed - Failed to validate book borrowed : ", err)
		return err
	}

	if count > 0 {
		r.Logger.Error("repo::ValidateBookBorrowed - Book already borrowed")
		return errors.New(constants.ErrBookAlreadyBorrowed)
	}

	return nil
}

func (r *BookBorrowedRepository) UpdateBookReturned(ctx context.Context, tx *sql.Tx, returnedDate time.Time, bookID, userID string) error {
	_, err := tx.ExecContext(ctx, r.DB.Rebind(queryUpdateBookReturned), returnedDate, bookID, userID)
	if err != nil {
		r.Logger.Error("repo::UpdateBookReturned - Failed to update book returned : ", err)
		return err
	}

	return nil
}

func (r *BookBorrowedRepository) ValidateBookReturned(ctx context.Context, tx *sql.Tx, bookID, userID string) error {
	var (
		count int
	)

	err := tx.QueryRowContext(ctx, r.DB.Rebind(queryValidateBookReturned), bookID, userID).Scan(&count)
	if err != nil {
		r.Logger.Error("repo::ValidateBookReturned - Failed to validate book returned : ", err)
		return err
	}

	if count > 0 {
		r.Logger.Error("repo::ValidateBookReturned - Book not borrowed")
		return errors.New(constants.ErrBookAlreadyReturned)
	}

	return nil
}
