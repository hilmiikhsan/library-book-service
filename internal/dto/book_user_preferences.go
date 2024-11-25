package dto

type CreateBookUserPreferencesRequest struct {
	UserID            string `json:"user_id"`
	PreferredCategory string `json:"preferred_category" validate:"required"`
}
