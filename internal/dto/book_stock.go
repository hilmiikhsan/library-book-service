package dto

type CreateBookStockRequest struct {
	BookID         string `json:"book_id" validate:"required"`
	TotalStock     int    `json:"total_stock" validate:"required,numeric"`
	AvailableStock int    `json:"available_stock" validate:"required,numeric"`
}

type UpdateBookStockRequest struct {
	ID             string `json:"id" validate:"required"`
	BookID         string `json:"book_id" validate:"required"`
	TotalStock     int    `json:"total_stock" validate:"required,numeric"`
	AvailableStock int    `json:"available_stock" validate:"required,numeric"`
}

type GetDetailBookStockResponse struct {
	ID             string     `json:"id"`
	Book           DetailBook `json:"book"`
	TotalStock     int        `json:"total_stock"`
	AvailableStock int        `json:"available_stock"`
}

type GetListBookStockResponse struct {
	BookStockList []BookStock `json:"book_stock_list"`
	Pagination    Pagination  `json:"pagination"`
}

type BookStock struct {
	ID             string     `json:"id"`
	Book           DetailBook `json:"book"`
	TotalStock     int        `json:"total_stock"`
	AvailableStock int        `json:"available_stock"`
}
