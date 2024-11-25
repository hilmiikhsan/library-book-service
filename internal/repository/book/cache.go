package book

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hilmiikhsan/library-book-service/helpers"
)

func (r *BookRepository) generateSearchBooksCacheKey(title *string, categoryID *string, authorID *string, limit, offset int) string {
	return fmt.Sprintf("search_books:%s:%s:%s:%d:%d",
		helpers.SafeString(title),
		helpers.SafeString(categoryID),
		helpers.SafeString(authorID),
		limit, offset,
	)
}

func (r *BookRepository) getCache(ctx context.Context, key string) ([]byte, error) {
	data, err := r.Redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

func (r *BookRepository) setCache(ctx context.Context, key string, data []byte, ttl int) error {
	return r.Redis.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
}
