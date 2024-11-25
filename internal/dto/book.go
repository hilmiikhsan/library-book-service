package dto

type CreateBookRequest struct {
	Title         string `json:"title" validate:"required,min=2,max=255"`
	AuthorID      string `json:"author_id" validate:"required"`
	CategoryID    string `json:"category_id" validate:"required"`
	Isbn          string `json:"isbn" validate:"required"`
	Description   string `json:"description" validate:"required"`
	PublishedDate string `json:"published_date" validate:"required"`
}

type UpdateBookRequest struct {
	ID            string `json:"id" validate:"required"`
	Title         string `json:"title" validate:"required,min=2,max=255"`
	AuthorID      string `json:"author_id" validate:"required"`
	CategoryID    string `json:"category_id" validate:"required"`
	Isbn          string `json:"isbn"`
	Description   string `json:"description"`
	PublishedDate string `json:"published_date" validate:"required"`
}

type GetDetailBookResponse struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Author        Author   `json:"author"`
	Category      Category `json:"category"`
	Description   string   `json:"description"`
	Isbn          string   `json:"isbn"`
	Stock         int      `json:"stock"`
	PublishedDate string   `json:"published_date"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

type GetListBookResponse struct {
	BookList   []Book     `json:"book_list"`
	Pagination Pagination `json:"pagination"`
}

type Book struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Isbn          string `json:"isbn"`
	PublishedDate string `json:"published_date"`
}

type DetailBook struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type SearchBookRequest struct {
	Title      string `json:"title,omitempty"`
	AuthorID   string `json:"author_id,omitempty"`
	CategoryID string `json:"category_id,omitempty"`
	Isbn       string `json:"isbn,omitempty"`
	Page       int    `json:"page,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

type GetListRecommendationsResponse struct {
	RecommendationList []Recommendations `json:"recommendation_list"`
	Pagination         Pagination        `json:"pagination"`
}

type Recommendations struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	AuthorID      string `json:"author_id"`
	CategoryID    string `json:"category_id"`
	Description   string `json:"description"`
	PublishedDate string `json:"published_date"`
}
