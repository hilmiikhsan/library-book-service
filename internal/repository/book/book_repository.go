package book

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/helpers"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type BookRepository struct {
	DB     *sqlx.DB
	Logger *logrus.Logger
	Redis  *redis.Client
}

func (r *BookRepository) InsertNewBook(ctx context.Context, book *models.Book) error {
	_, err := r.DB.ExecContext(ctx, r.DB.Rebind(queryInsertNewBook),
		book.Title,
		book.AuthorID,
		book.CategoryID,
		book.Isbn,
		book.Description,
		book.PublishedDate,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if !ok {
			r.Logger.Error("repo::InsertNewBook - Failed to insert new book : ", err)
			return err
		}

		switch pqErr.Code.Name() {
		case "unique_violation":
			r.Logger.Error("repo::InsertNewBook - isbn already exist: ", err)
			return errors.New(constants.ErrIsbnAlreadyExist)
		default:
			r.Logger.Error("repo::InsertNewBook - Failed to insert new book : ", err)
			return err
		}
	}

	return nil
}

func (r *BookRepository) FindBookByID(ctx context.Context, id string) (*models.Book, error) {
	var (
		res      = new(models.Book)
		cacheKey = fmt.Sprintf("book:%s", id)
	)

	cachedData, err := r.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cachedData), res)
		if err == nil {
			r.Logger.Info("category::FindBookByID - Data retrieved from cache")
			return res, nil
		}
		r.Logger.Warn("category::FindBookByID - Failed to unmarshal cache data: ", err)
	}

	err = r.DB.GetContext(ctx, res, r.DB.Rebind(queryFindBookByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Error("repo::FindBookByID - Book doesnt exist")
			return res, errors.New(constants.ErrBookNotFound)
		}

		r.Logger.Error("repo::FindBookByID - failed to find book by id: ", err)
		return nil, err
	}

	dataToCache, err := json.Marshal(res)
	if err != nil {
		r.Logger.Warn("category::FindBookByID - Failed to marshal data for caching: ", err)
	} else {
		err = r.Redis.Set(ctx, cacheKey, dataToCache, 5*time.Minute).Err()
		if err != nil {
			r.Logger.Warn("category::FindBookByID - Failed to cache data: ", err)
		}
	}

	return res, nil
}

func (r *BookRepository) FindAllBook(ctx context.Context, limit, offset int) ([]models.Book, error) {
	var (
		res      = make([]models.Book, 0)
		cacheKey = fmt.Sprintf("books:limit:%d:offset:%d", limit, offset)
	)

	cachedData, err := r.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cachedData), &res)
		if err == nil {
			r.Logger.Info("category::FindAllBook - Data retrieved from cache")
			return res, nil
		}
		r.Logger.Warn("category::FindAllBook - Failed to unmarshal cache data: ", err)
	}

	err = r.DB.SelectContext(ctx, &res, r.DB.Rebind(queryFindAllBook), limit, offset)
	if err != nil {
		r.Logger.Error("repo::FindAllBook - failed to find all book: ", err)
		return nil, err
	}

	dataToCache, err := json.Marshal(res)
	if err != nil {
		r.Logger.Warn("category::FindAllBook - Failed to marshal data for caching: ", err)
	} else {
		err = r.Redis.Set(ctx, cacheKey, dataToCache, 5*time.Minute).Err()
		if err != nil {
			r.Logger.Warn("category::FindAllBook - Failed to cache data: ", err)
		}
	}

	return res, nil
}

func (r *BookRepository) UpdateNewBook(ctx context.Context, book *models.Book) error {
	query := `
		UPDATE books
		SET
			title = ?,
			author_id = ?,
			category_id = ?,
			description = ?,
			published_date = ?,
			update_at = now()
	`

	args := []interface{}{
		book.Title,
		book.AuthorID,
		book.CategoryID,
		book.Description,
		book.PublishedDate,
	}

	if book.Isbn != nil {
		query += ", isbn = ?"
		args = append(args, *book.Isbn)
	}

	query += " WHERE id = ?"
	args = append(args, book.ID)

	_, err := r.DB.ExecContext(ctx, r.DB.Rebind(query), args...)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if !ok {
			r.Logger.Error("repo::UpdateNewBook - Failed to insert new book : ", err)
			return err
		}

		switch pqErr.Code.Name() {
		case "unique_violation":
			r.Logger.Error("repo::UpdateNewBook - isbn already exist: ", err)
			return errors.New(constants.ErrIsbnAlreadyExist)
		default:
			r.Logger.Error("repo::UpdateNewBook - Failed to insert new book : ", err)
			return err
		}
	}

	return nil
}

func (r *BookRepository) DeleteBookByID(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, r.DB.Rebind(queryDeleteBookByID), id)
	if err != nil {
		r.Logger.Error("repo::DeleteBookByID - failed to delete book by id: ", err)
		return err
	}

	return nil
}

func (r *BookRepository) SearchBooks(ctx context.Context, title *string, categoryID *string, authorID *string, limit, offset int) ([]models.Book, error) {
	var res []models.Book
	cacheKey := r.generateSearchBooksCacheKey(title, categoryID, authorID, limit, offset)

	cachedData, err := r.getCache(ctx, cacheKey)
	if err == nil && cachedData != nil {
		if err := helpers.UnmarshalJSON(cachedData, &res); err == nil {
			r.Logger.Info("repo::SearchBooks - returned data from cache")
			return res, nil
		}
	}

	query := `
        SELECT id, title, author_id, category_id, isbn, description, published_date, created_at, updated_at
        FROM books
        WHERE TRUE
    `

	args := []interface{}{}

	if title != nil {
		query += " AND title ILIKE ?"
		args = append(args, "%"+*title+"%")
	}
	if categoryID != nil && *categoryID != "" {
		query += " AND category_id = ?"
		args = append(args, *categoryID)
	}
	if authorID != nil && *authorID != "" {
		query += " AND author_id = ?"
		args = append(args, *authorID)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err = r.DB.SelectContext(ctx, &res, r.DB.Rebind(query), args...)
	if err != nil {
		r.Logger.Error("repo::SearchBooks - failed to search books: ", err)
		return nil, err
	}

	if err := r.setCache(ctx, cacheKey, helpers.MarshalJSON(res), 300); err != nil {
		r.Logger.Warn("repo::SearchBooks - failed to set cache: ", err)
	}

	return res, nil
}

func (r *BookRepository) GetRecommendations(ctx context.Context, userID string, limit, offset int) ([]models.Book, error) {
	var (
		books    []models.Book
		cacheKey = fmt.Sprintf("recommendations:%s:%d:%d", userID, limit, offset)
	)

	cachedData, err := r.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cachedData), &books)
		if err == nil {
			r.Logger.Info("category::GetRecommendations - Data retrieved from cache")
			return books, nil
		}
		r.Logger.Warn("category::GetRecommendations - Failed to unmarshal cache data: ", err)
	}

	err = r.DB.SelectContext(ctx, &books, r.DB.Rebind(queryGetRecommendations), userID, limit, offset)
	if err != nil {
		r.Logger.Error("repo::GetRecommendations - Failed to fetch recommendations: ", err)
		return nil, err
	}

	dataToCache, err := json.Marshal(books)
	if err != nil {
		r.Logger.Warn("category::GetRecommendations - Failed to marshal data for caching: ", err)
	} else {
		err = r.Redis.Set(ctx, cacheKey, dataToCache, 5*time.Minute).Err()
		if err != nil {
			r.Logger.Warn("category::GetRecommendations - Failed to cache data: ", err)
		}
	}

	return books, nil
}
