package interfaces

import (
	"context"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/models"
)

type IBookBorrowedRepository interface {
	InsertNewBookBorrowed(ctx context.Context, tx *sql.Tx, bookBorrowed *models.BookBorrowed) error
	ValidateBookBorrowed(ctx context.Context, tx *sql.Tx, bookID, userID string) error
	UpdateBookReturned(ctx context.Context, tx *sql.Tx, returnedDate time.Time, bookID, userID string) error
	ValidateBookReturned(ctx context.Context, tx *sql.Tx, bookID, userID string) error
}

type IBookBorrowedService interface {
	BookBorrowed(ctx context.Context, req *dto.BookBorrowedRequest, userID string) error
	BookReturned(ctx context.Context, req *dto.BookReturnedRequest, userID string) error
}

type IBookBorrowedHandler interface {
	BookBorrowed(*gin.Context)
	BookReturned(*gin.Context)
}
