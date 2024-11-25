package interfaces

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/models"
)

type IBookRepository interface {
	InsertNewBook(ctx context.Context, book *models.Book) error
	FindBookByID(ctx context.Context, id string) (*models.Book, error)
	FindAllBook(ctx context.Context, limit, offset int) ([]models.Book, error)
	UpdateNewBook(ctx context.Context, book *models.Book) error
	DeleteBookByID(ctx context.Context, id string) error
	SearchBooks(ctx context.Context, title *string, categoryID *string, authorID *string, limit, offset int) ([]models.Book, error)
	GetRecommendations(ctx context.Context, userID string, limit, offset int) ([]models.Book, error)
}

type IBookService interface {
	CreateBook(ctx context.Context, req *dto.CreateBookRequest) error
	GetDetailBook(ctx context.Context, id string) (*dto.GetDetailBookResponse, error)
	GetListBook(ctx context.Context, limit, offset int) (*dto.GetListBookResponse, error)
	UpdateBook(ctx context.Context, req *dto.UpdateBookRequest) error
	DeleteBook(ctx context.Context, id string) error
	SearchBooks(ctx context.Context, req *dto.SearchBookRequest) (*dto.GetListBookResponse, error)
	GetRecommendations(ctx context.Context, userID string, limit, offset int) (*dto.GetListRecommendationsResponse, error)
}

type IBookHandler interface {
	CreateBook(*gin.Context)
	GetDetailBook(*gin.Context)
	GetListBook(*gin.Context)
	UpdateBook(*gin.Context)
	DeleteBook(*gin.Context)
	SearchBooks(*gin.Context)
	GetRecommendations(*gin.Context)
}
