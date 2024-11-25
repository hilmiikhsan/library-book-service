package models

import "time"

type BookUserPreferences struct {
	ID                string    `db:"id"`
	UserID            string    `db:"user_id"`
	PreferredCategory string    `db:"preferred_category"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}
