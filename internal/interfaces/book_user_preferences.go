package interfaces

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/models"
)

type IBookUserPreferencesRepository interface {
	InsertNewBookUserPreferences(ctx context.Context, bookBorrowed *models.BookUserPreferences) error
}

type IBookUserPreferencesService interface {
	CreateBookUserPreferences(ctx context.Context, req *dto.CreateBookUserPreferencesRequest) error
}

type IBookUserPreferencesHandler interface {
	CreateBookUserPreferences(*gin.Context)
}
