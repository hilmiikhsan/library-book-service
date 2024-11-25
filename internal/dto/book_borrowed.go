package dto

type BookBorrowedRequest struct {
	BookID  string `json:"book_id" validate:"required"`
	DueDate string `json:"due_date" validate:"required"`
}

type UpdateBookBorrowedRequest struct {
	ID      string `json:"id" validate:"required"`
	BookID  string `json:"book_id" validate:"required"`
	DueDate string `json:"due_date" validate:"required"`
}

type BookReturnedRequest struct {
	BookID       string `json:"book_id" validate:"required"`
	ReturnedDate string `json:"returned_date" validate:"required"`
}

// type GetDetailBookStockResponse struct {
// 	ID             string     `json:"id"`
// 	Book           DetailBook `json:"book"`
// 	TotalStock     int        `json:"total_stock"`
// 	AvailableStock int        `json:"available_stock"`
// }

// type GetListBookStockResponse struct {
// 	BookStockList []BookStock `json:"book_stock_list"`
// 	Pagination    Pagination  `json:"pagination"`
// }

// type BookStock struct {
// 	ID             string     `json:"id"`
// 	Book           DetailBook `json:"book"`
// 	TotalStock     int        `json:"total_stock"`
// 	AvailableStock int        `json:"available_stock"`
// }
