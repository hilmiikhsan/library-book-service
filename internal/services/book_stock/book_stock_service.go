package book_stock

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/interfaces"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/sirupsen/logrus"
)

type BookStockService struct {
	BookStockRepo interfaces.IBookStockRepository
	BookRepo      interfaces.IBookRepository
	Logger        *logrus.Logger
}

func (s *BookStockService) CreateBookStock(ctx context.Context, req *dto.CreateBookStockRequest) error {
	bookID, _ := uuid.Parse(req.BookID)

	_, err := s.BookRepo.FindBookByID(ctx, req.BookID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrBookNotFound) {
			s.Logger.Error("service::CreateBookStock - book not found")
			return err
		}

		s.Logger.Error("service::CreateBookStock - failed to get detail book: ", err)
		return err
	}

	countData, err := s.BookStockRepo.ValidateBookStockByBookID(ctx, req.BookID)
	if err != nil {
		s.Logger.Error("service::CreateBookStock - failed to validate BookStock by book id: ", err)
		return err
	}

	if countData > 0 {
		s.Logger.Error("service::CreateBookStock - BookStock already exist")
		return errors.New(constants.ErrBookStockAlreadyExist)
	}

	err = s.BookStockRepo.InsertNewBookStock(ctx, &models.BookStock{
		BookID:         bookID,
		TotalStock:     req.TotalStock,
		AvailableStock: req.AvailableStock,
	})
	if err != nil {
		s.Logger.Error("service::CreateBookStock - failed to insert new BookStock: ", err)
		return err
	}

	return nil
}

func (s *BookStockService) GetDetailBookStock(ctx context.Context, id string) (*dto.GetDetailBookStockResponse, error) {
	bookStockData, err := s.BookStockRepo.FindBookStockByID(ctx, id)
	if err != nil {
		s.Logger.Error("service::GetBookStockDetail - failed to find BookStock by id: ", err)
		return &dto.GetDetailBookStockResponse{}, err
	}

	return &dto.GetDetailBookStockResponse{
		ID: bookStockData.ID.String(),
		Book: dto.DetailBook{
			ID:    bookStockData.BookID.String(),
			Title: bookStockData.BookTitle,
		},
		TotalStock:     bookStockData.TotalStock,
		AvailableStock: bookStockData.AvailableStock,
	}, nil
}

func (s *BookStockService) GetListBookStock(ctx context.Context, limit, offset int) (*dto.GetListBookStockResponse, error) {
	pageSize := limit
	pageIndex := (offset - 1) * limit

	bookStockData, err := s.BookStockRepo.FindAllBookStock(ctx, pageSize, pageIndex)
	if err != nil {
		s.Logger.Error("service::GetListBookStock - failed to find all BookStock: ", err)
		return nil, err
	}

	bookStocks := make([]dto.BookStock, 0)
	for _, bookStock := range bookStockData {
		bookStocks = append(bookStocks, dto.BookStock{
			ID: bookStock.ID.String(),
			Book: dto.DetailBook{
				ID:    bookStock.BookID.String(),
				Title: bookStock.BookTitle,
			},
			TotalStock:     bookStock.TotalStock,
			AvailableStock: bookStock.AvailableStock,
		})
	}

	pagination := dto.Pagination{
		Page:  offset,
		Limit: limit,
	}

	response := &dto.GetListBookStockResponse{
		BookStockList: bookStocks,
		Pagination:    pagination,
	}

	return response, nil
}

func (s *BookStockService) UpdateBookStock(ctx context.Context, req *dto.UpdateBookStockRequest) error {
	bookStockData, err := s.BookStockRepo.FindBookStockByID(ctx, req.ID)
	if err != nil {
		s.Logger.Error("service::UpdateBookStock - failed to find BookStock by id: ", err)
		return err
	}

	if len(bookStockData.ID) == 0 {
		s.Logger.Error("service::UpdateBookStock - BookStock not found")
		return errors.New(constants.ErrBookStockNotFound)
	}

	bookID, _ := uuid.Parse(req.BookID)

	mappingBookStockData := &models.BookStock{
		ID:             bookStockData.ID,
		BookID:         bookID,
		TotalStock:     req.TotalStock,
		AvailableStock: req.AvailableStock,
	}

	err = s.BookStockRepo.UpdateNewBookStock(ctx, mappingBookStockData)
	if err != nil {
		s.Logger.Error("service::UpdateBookStock - failed to update BookStock: ", err)
		return err
	}

	return nil
}

func (s *BookStockService) DeleteBookStock(ctx context.Context, id string) error {
	bookStockData, err := s.BookStockRepo.FindBookStockByID(ctx, id)
	if err != nil {
		s.Logger.Error("service::DeleteBookStock - failed to find BookStock by id: ", err)
		return err
	}

	if len(bookStockData.ID) == 0 {
		s.Logger.Error("service::DeleteBookStock - BookStock not found")
		return errors.New(constants.ErrBookStockNotFound)
	}

	err = s.BookStockRepo.DeleteBookStockByID(ctx, bookStockData.ID.String())
	if err != nil {
		s.Logger.Error("service::DeleteBookStock - failed to delete BookStock: ", err)
		return err
	}

	return nil
}
