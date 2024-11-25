package interfaces

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/models"
)

type IBookStockRepository interface {
	InsertNewBookStock(ctx context.Context, bookStock *models.BookStock) error
	FindBookStockByID(ctx context.Context, id string) (*models.BookStock, error)
	FindAllBookStock(ctx context.Context, limit, offset int) ([]models.BookStock, error)
	UpdateNewBookStock(ctx context.Context, bookStock *models.BookStock) error
	DeleteBookStockByID(ctx context.Context, id string) error
	ValidateBookStockByBookID(ctx context.Context, bookID string) (int, error)
	DecrementAvailableStock(ctx context.Context, tx *sql.Tx, bookID string, stock int) error
	LockBookStock(ctx context.Context, tx *sql.Tx, bookID string) error
	IncrementAvailableStock(ctx context.Context, tx *sql.Tx, bookID string, stock int) error
	LockBookStockReturned(ctx context.Context, tx *sql.Tx, bookID string) error
}

type IBookStockService interface {
	CreateBookStock(ctx context.Context, req *dto.CreateBookStockRequest) error
	GetDetailBookStock(ctx context.Context, id string) (*dto.GetDetailBookStockResponse, error)
	GetListBookStock(ctx context.Context, limit, offset int) (*dto.GetListBookStockResponse, error)
	UpdateBookStock(ctx context.Context, req *dto.UpdateBookStockRequest) error
	DeleteBookStock(ctx context.Context, id string) error
}

type IBookStockHandler interface {
	CreateBookStock(*gin.Context)
	GetDetailBookStock(*gin.Context)
	GetListBookStock(*gin.Context)
	UpdateBookStock(*gin.Context)
	DeleteBookStock(*gin.Context)
}
