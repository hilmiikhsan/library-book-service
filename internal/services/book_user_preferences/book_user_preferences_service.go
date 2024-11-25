package book_user_preferences

import (
	"context"
	"strings"

	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/interfaces"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/sirupsen/logrus"
)

type BookUserPreferencesService struct {
	BookUserPreferencesRepo interfaces.IBookUserPreferencesRepository
	External                interfaces.IExternal
	Logger                  *logrus.Logger
}

func (s *BookUserPreferencesService) CreateBookUserPreferences(ctx context.Context, req *dto.CreateBookUserPreferencesRequest) error {
	_, err := s.External.GetDetailCategory(ctx, req.PreferredCategory)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrCategoryNotFound) {
			s.Logger.Error("service::CreateBookUserPreferences - category not found")
			return err
		}

		s.Logger.Error("service::CreateBookUserPreferences - failed to get detail category: ", err)
		return err
	}

	err = s.BookUserPreferencesRepo.InsertNewBookUserPreferences(ctx, &models.BookUserPreferences{
		UserID:            req.UserID,
		PreferredCategory: req.PreferredCategory,
	})
	if err != nil {
		s.Logger.Error("service::CreateBookUserPreferences - failed to insert new BookUserPreferences: ", err)
		return err
	}

	return nil
}
