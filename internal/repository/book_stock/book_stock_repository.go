package BookStock

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type BookStockRepository struct {
	DB     *sqlx.DB
	Logger *logrus.Logger
	Redis  *redis.Client
}

func (r *BookStockRepository) InsertNewBookStock(ctx context.Context, bookStock *models.BookStock) error {
	_, err := r.DB.ExecContext(ctx, r.DB.Rebind(queryInsertNewBookStock),
		bookStock.BookID,
		bookStock.TotalStock,
		bookStock.AvailableStock,
	)
	if err != nil {
		r.Logger.Error("repo::InsertNewBookStock - Failed to insert new book stock : ", err)
		return err
	}

	return nil
}

func (r *BookStockRepository) FindBookStockByID(ctx context.Context, id string) (*models.BookStock, error) {
	var (
		res      = new(models.BookStock)
		cacheKey = fmt.Sprintf("book_stock:%s", id)
	)

	cachedData, err := r.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cachedData), res)
		if err == nil {
			r.Logger.Info("category::FindBookStockByID - Data retrieved from cache")
			return res, nil
		}
		r.Logger.Warn("category::FindBookStockByID - Failed to unmarshal cache data: ", err)
	}

	err = r.DB.GetContext(ctx, res, r.DB.Rebind(queryFindBookStockByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Error("repo::FindBookStockByID - BookStock doesnt exist")
			return res, errors.New(constants.ErrBookStockNotFound)
		}

		r.Logger.Error("repo::FindBookStockByID - failed to find book stock by id: ", err)
		return nil, err
	}

	dataToCache, err := json.Marshal(res)
	if err != nil {
		r.Logger.Warn("category::FindBookStockByID - Failed to marshal data for caching: ", err)
	} else {
		err = r.Redis.Set(ctx, cacheKey, dataToCache, 5*time.Minute).Err()
		if err != nil {
			r.Logger.Warn("category::FindBookStockByID - Failed to cache data: ", err)
		}
	}

	return res, nil
}

func (r *BookStockRepository) FindAllBookStock(ctx context.Context, limit, offset int) ([]models.BookStock, error) {
	var (
		res      = make([]models.BookStock, 0)
		cacheKey = fmt.Sprintf("book_stock:limit:%d:offset:%d", limit, offset)
	)

	cachedData, err := r.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cachedData), &res)
		if err == nil {
			r.Logger.Info("category::FindAllBookStock - Data retrieved from cache")
			return res, nil
		}
		r.Logger.Warn("category::FindAllBookStock - Failed to unmarshal cache data: ", err)
	}

	err = r.DB.SelectContext(ctx, &res, r.DB.Rebind(queryFindAllBookStock), limit, offset)
	if err != nil {
		r.Logger.Error("repo::FindAllBookStock - failed to find all book stock: ", err)
		return nil, err
	}

	dataToCache, err := json.Marshal(res)
	if err != nil {
		r.Logger.Warn("category::FindAllBookStock - Failed to marshal data for caching: ", err)
	} else {
		err = r.Redis.Set(ctx, cacheKey, dataToCache, 5*time.Minute).Err()
		if err != nil {
			r.Logger.Warn("category::FindAllBookStock - Failed to cache data: ", err)
		}
	}

	return res, nil
}

func (r *BookStockRepository) UpdateNewBookStock(ctx context.Context, bookStock *models.BookStock) error {
	_, err := r.DB.ExecContext(ctx, r.DB.Rebind(queryUpdateBookStock),
		bookStock.BookID,
		bookStock.TotalStock,
		bookStock.AvailableStock,
		bookStock.ID,
	)
	if err != nil {
		r.Logger.Error("repo::UpdateNewBookStock - failed to update book stock: ", err)
		return err
	}

	return nil
}

func (r *BookStockRepository) DeleteBookStockByID(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, r.DB.Rebind(queryDeleteBookStockByID), id)
	if err != nil {
		r.Logger.Error("repo::DeleteBookStockByID - failed to delete BookStock by id: ", err)
		return err
	}

	return nil
}

func (r *BookStockRepository) ValidateBookStockByBookID(ctx context.Context, bookID string) (int, error) {
	var count int

	err := r.DB.GetContext(ctx, &count, r.DB.Rebind(queryCountBookByBookID), bookID)
	if err != nil {
		r.Logger.Error("repo::ValidateBookStockByBookID - failed to count book by book id: ", err)
		return 0, err
	}

	return count, nil
}

func (r *BookStockRepository) DecrementAvailableStock(ctx context.Context, tx *sql.Tx, bookID string, stock int) error {
	result, err := tx.ExecContext(ctx, r.DB.Rebind(queryDecrementAvailableStock), stock, bookID, stock)
	if err != nil {
		r.Logger.Error("repo::UpdateAvailableStock - failed to update available stock: ", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.Logger.Error("repo::UpdateAvailableStock - failed to get rows affected: ", err)
		return err
	}

	if rowsAffected == 0 {
		r.Logger.Error("repo::UpdateAvailableStock - insufficient stock")
		return errors.New(constants.ErrInsufficientStock)
	}

	return nil
}

func (r *BookStockRepository) LockBookStock(ctx context.Context, tx *sql.Tx, bookID string) error {
	_, err := tx.ExecContext(ctx, r.DB.Rebind(queryLockBookStock), bookID)
	if err != nil {
		r.Logger.Error("repo::LockBookStock - failed to lock book stock: ", err)
		return err
	}

	return nil
}

func (r *BookStockRepository) IncrementAvailableStock(ctx context.Context, tx *sql.Tx, bookID string, stock int) error {
	_, err := tx.ExecContext(ctx, r.DB.Rebind(queryIncrementAvailableStock), stock, bookID)
	if err != nil {
		r.Logger.Error("repo::IncrementAvailableStock - failed to increment available stock: ", err)
		return err
	}

	return nil
}

func (r *BookStockRepository) LockBookStockReturned(ctx context.Context, tx *sql.Tx, bookID string) error {
	_, err := tx.ExecContext(ctx, r.DB.Rebind(queryLockBookStockReturned), bookID)
	if err != nil {
		r.Logger.Error("repo::LockBookStockReturned - failed to lock book stock returned: ", err)
		return err
	}

	return nil
}
