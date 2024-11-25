package book_borrowed

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/helpers"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/interfaces"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type BookBorrowedService struct {
	BookBorrowedRepo interfaces.IBookBorrowedRepository
	BookStockRepo    interfaces.IBookStockRepository
	Logger           *logrus.Logger
	DB               *sqlx.DB
}

func (s *BookBorrowedService) BookBorrowed(ctx context.Context, req *dto.BookBorrowedRequest, userID string) error {
	userId, _ := uuid.Parse(userID)
	bookId, _ := uuid.Parse(req.BookID)

	dueDate, err := helpers.ParseDate(req.DueDate, constants.DateTimeFormat)
	if err != nil {
		s.Logger.Error("service::BookBorrowed - failed to parse due date: ", err)
		return errors.New(constants.ErrInvalidFormatDate)
	}

	countData, err := s.BookStockRepo.ValidateBookStockByBookID(ctx, req.BookID)
	if err != nil {
		s.Logger.Error("service::BookBorrowed - failed to validate book stock: ", err)
		return err
	}

	if countData <= 0 {
		s.Logger.Error("service::BookBorrowed - book stock not found")
		return errors.New(constants.ErrBookStockNotFound)
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		s.Logger.Error("service::BookBorrowed - failed to begin transaction: ", err)
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				s.Logger.Error("service::BookBorrowed - failed to rollback transaction: ", rollbackErr)
			}
		}
	}()

	err = s.BookBorrowedRepo.ValidateBookBorrowed(ctx, tx, req.BookID, userID)
	if err != nil {
		s.Logger.Error("service::BookBorrowed - failed to validate book borrowed: ", err)
		return err
	}

	err = s.BookStockRepo.LockBookStock(ctx, tx, req.BookID)
	if err != nil {
		s.Logger.Error("service::BookBorrowed - failed to lock book stock: ", err)
		return err
	}

	err = s.BookBorrowedRepo.InsertNewBookBorrowed(ctx, tx, &models.BookBorrowed{
		UserID:  userId,
		BookID:  bookId,
		DueDate: dueDate,
	})
	if err != nil {
		s.Logger.Error("service::BookBorrowed - failed to insert new book borrowed: ", err)
		return err
	}

	err = s.BookStockRepo.DecrementAvailableStock(ctx, tx, req.BookID, 1)
	if err != nil {
		s.Logger.Error("service::BookBorrowed - failed to update available stock: ", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		s.Logger.Error("service::BookBorrowed - failed to commit transaction: ", err)
		return err
	}

	return nil
}

func (s *BookBorrowedService) BookReturned(ctx context.Context, req *dto.BookReturnedRequest, userID string) error {
	returnedDate, err := helpers.ParseDate(req.ReturnedDate, constants.DateTimeFormat)
	if err != nil {
		s.Logger.Error("service::BookBorrowed - failed to parse returned date: ", err)
		return errors.New(constants.ErrInvalidFormatDate)
	}

	countData, err := s.BookStockRepo.ValidateBookStockByBookID(ctx, req.BookID)
	if err != nil {
		s.Logger.Error("service::BookBorrowed - failed to validate book stock: ", err)
		return err
	}

	if countData <= 0 {
		s.Logger.Error("service::BookBorrowed - book stock not found")
		return errors.New(constants.ErrBookStockNotFound)
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		s.Logger.Error("service::BookReturned - failed to begin transaction: ", err)
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				s.Logger.Error("service::BookReturned - failed to rollback transaction: ", rollbackErr)
			}
		}
	}()

	err = s.BookStockRepo.LockBookStockReturned(ctx, tx, req.BookID)
	if err != nil {
		s.Logger.Error("service::BookReturned - failed to lock book stock: ", err)
		return err
	}

	err = s.BookBorrowedRepo.ValidateBookReturned(ctx, tx, req.BookID, userID)
	if err != nil {
		s.Logger.Error("service::BookReturned - failed to validate book returned: ", err)
		return err
	}

	err = s.BookBorrowedRepo.UpdateBookReturned(ctx, tx, returnedDate, req.BookID, userID)
	if err != nil {
		s.Logger.Error("service::BookReturned - failed to update book returned: ", err)
		return err
	}

	err = s.BookStockRepo.IncrementAvailableStock(ctx, tx, req.BookID, 1)
	if err != nil {
		s.Logger.Error("service::BookReturned - failed to update available stock: ", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		s.Logger.Error("service::BookReturned - failed to commit transaction: ", err)
		return err
	}

	return nil
}
