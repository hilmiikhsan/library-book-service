package book_user_preferences

import (
	"context"

	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type BookUserPreferencesRepository struct {
	DB     *sqlx.DB
	Logger *logrus.Logger
}

func (r *BookUserPreferencesRepository) InsertNewBookUserPreferences(ctx context.Context, bookBorrowed *models.BookUserPreferences) error {
	var count int
	err := r.DB.GetContext(ctx, &count, r.DB.Rebind(queryCountUserPreferences), bookBorrowed.UserID, bookBorrowed.PreferredCategory)
	if err != nil {
		r.Logger.Error("repo::InsertNewBookUserPreferences - Failed to check existing preferences: ", err)
		return err
	}

	if count > 0 {
		_, err := r.DB.ExecContext(ctx, r.DB.Rebind(queryDeleteExistingUserPreferences), bookBorrowed.UserID, bookBorrowed.PreferredCategory)
		if err != nil {
			r.Logger.Error("repo::InsertNewBookUserPreferences - Failed to delete existing preferences: ", err)
			return err
		}
	}

	_, err = r.DB.ExecContext(ctx, r.DB.Rebind(queryInsertNewBookUserPreferences), bookBorrowed.UserID, bookBorrowed.PreferredCategory)
	if err != nil {
		r.Logger.Error("repo::InsertNewBookUserPreferences - Failed to insert new book user preferences: ", err)
		return err
	}

	return nil
}
