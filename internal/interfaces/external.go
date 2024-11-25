package interfaces

import (
	"context"

	"github.com/hilmiikhsan/library-book-service/internal/models"
)

type IExternal interface {
	ValidateToken(ctx context.Context, token string) (models.TokenData, error)
	GetDetailAuthor(ctx context.Context, id string) (models.AuthorModel, error)
	GetDetailCategory(ctx context.Context, id string) (models.CategoryModel, error)
}
