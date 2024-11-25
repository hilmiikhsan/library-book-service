package book

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/helpers"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/interfaces"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/sirupsen/logrus"
)

type BookService struct {
	BookRepo interfaces.IBookRepository
	External interfaces.IExternal
	Logger   *logrus.Logger
}

func (s *BookService) CreateBook(ctx context.Context, req *dto.CreateBookRequest) error {
	authorID, _ := uuid.Parse(req.AuthorID)
	categoryID, _ := uuid.Parse(req.CategoryID)

	_, err := s.External.GetDetailAuthor(ctx, req.AuthorID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrAuthorNotFound) {
			s.Logger.Error("service::CreateBook - author not found")
			return err
		}

		s.Logger.Error("service::CreateBook - failed to get detail author: ", err)
		return err
	}

	_, err = s.External.GetDetailCategory(ctx, req.CategoryID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrCategoryNotFound) {
			s.Logger.Error("service::CreateBook - category not found")
			return err
		}

		s.Logger.Error("service::CreateBook - failed to get detail category: ", err)
		return err
	}

	publishedDate, err := helpers.ParseDate(req.PublishedDate, constants.DateTimeFormat)
	if err != nil {
		s.Logger.Error("service::CreateBook - failed to parse published date: ", err)
		return errors.New(constants.ErrInvalidFormatDate)
	}

	err = s.BookRepo.InsertNewBook(ctx, &models.Book{
		Title:         req.Title,
		AuthorID:      authorID,
		CategoryID:    categoryID,
		Isbn:          &req.Isbn,
		Description:   req.Description,
		PublishedDate: publishedDate,
	})
	if err != nil {
		s.Logger.Error("service::CreateBook - failed to insert new book: ", err)
		return err
	}

	return nil
}

func (s *BookService) GetDetailBook(ctx context.Context, id string) (*dto.GetDetailBookResponse, error) {
	bookData, err := s.BookRepo.FindBookByID(ctx, id)
	if err != nil {
		s.Logger.Error("service::GetBookDetail - failed to find book by id: ", err)
		return &dto.GetDetailBookResponse{}, err
	}

	authorData, err := s.External.GetDetailAuthor(ctx, bookData.AuthorID.String())
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrAuthorNotFound) {
			s.Logger.Error("service::GetDetailBook - author not found")
			return &dto.GetDetailBookResponse{}, err
		}

		s.Logger.Error("service::GetDetailBook - failed to get detail author: ", err)
		return &dto.GetDetailBookResponse{}, err
	}

	categoryData, err := s.External.GetDetailCategory(ctx, bookData.CategoryID.String())
	if err != nil {
		s.Logger.Error("service::GetDetailBook - failed to get detail category: ", err)
		return &dto.GetDetailBookResponse{}, err
	}

	return &dto.GetDetailBookResponse{
		ID:    bookData.ID.String(),
		Title: bookData.Title,
		Author: dto.Author{
			ID:   authorData.ID,
			Name: authorData.Name,
		},
		Category: dto.Category{
			ID:   categoryData.ID,
			Name: categoryData.Name,
		},
		Description:   bookData.Description,
		Isbn:          *bookData.Isbn,
		PublishedDate: bookData.PublishedDate.Format(constants.DateTimeFormat),
		CreatedAt:     bookData.CreatedAt.Format(constants.DateTimeFormat),
		UpdatedAt:     bookData.UpdatedAt.Format(constants.DateTimeFormat),
	}, nil
}

func (s *BookService) GetListBook(ctx context.Context, limit, offset int) (*dto.GetListBookResponse, error) {
	pageSize := limit
	pageIndex := (offset - 1) * limit

	bookData, err := s.BookRepo.FindAllBook(ctx, pageSize, pageIndex)
	if err != nil {
		s.Logger.Error("service::GetListBook - failed to find all book: ", err)
		return nil, err
	}

	books := make([]dto.Book, 0)
	for _, book := range bookData {
		books = append(books, dto.Book{
			ID:            book.ID.String(),
			Title:         book.Title,
			Description:   book.Description,
			Isbn:          *book.Isbn,
			PublishedDate: book.PublishedDate.Format(constants.DateTimeFormat),
		})
	}

	pagination := dto.Pagination{
		Page:  offset,
		Limit: limit,
	}

	response := &dto.GetListBookResponse{
		BookList:   books,
		Pagination: pagination,
	}

	return response, nil
}

func (s *BookService) UpdateBook(ctx context.Context, req *dto.UpdateBookRequest) error {
	bookData, err := s.BookRepo.FindBookByID(ctx, req.ID)
	if err != nil {
		s.Logger.Error("service::UpdateBook - failed to find book by id: ", err)
		return err
	}

	if len(bookData.ID) == 0 {
		s.Logger.Error("service::UpdateBook - Book not found")
		return errors.New(constants.ErrBookNotFound)
	}

	publishedDate, err := helpers.ParseDate(req.PublishedDate, constants.DateTimeFormat)
	if err != nil {
		s.Logger.Error("service::UpdateBook - failed to parse published date: ", err)
		return errors.New(constants.ErrInvalidFormatDate)
	}

	_, err = s.External.GetDetailAuthor(ctx, req.AuthorID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrAuthorNotFound) {
			s.Logger.Error("service::UpdateBook - author not found")
			return err
		}

		s.Logger.Error("service::UpdateBook - failed to get detail author: ", err)
		return err
	}

	_, err = s.External.GetDetailCategory(ctx, req.CategoryID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrCategoryNotFound) {
			s.Logger.Error("service::UpdateBook - category not found")
			return err
		}

		s.Logger.Error("service::UpdateBook - failed to get detail category: ", err)
		return err
	}

	authorID, _ := uuid.Parse(req.AuthorID)
	categoryID, _ := uuid.Parse(req.CategoryID)

	mappingBookData := &models.Book{
		ID:            bookData.ID,
		Title:         req.Title,
		AuthorID:      authorID,
		CategoryID:    categoryID,
		Description:   req.Description,
		PublishedDate: publishedDate,
	}

	if req.Isbn != "" {
		mappingBookData.Isbn = &req.Isbn
	}

	err = s.BookRepo.UpdateNewBook(ctx, mappingBookData)
	if err != nil {
		s.Logger.Error("service::UpdateBook - failed to update book: ", err)
		return err
	}

	return nil
}

func (s *BookService) DeleteBook(ctx context.Context, id string) error {
	bookData, err := s.BookRepo.FindBookByID(ctx, id)
	if err != nil {
		s.Logger.Error("service::DeleteBook - failed to find book by id: ", err)
		return err
	}

	if len(bookData.ID) == 0 {
		s.Logger.Error("service::DeleteBook - Book not found")
		return errors.New(constants.ErrBookNotFound)
	}

	err = s.BookRepo.DeleteBookByID(ctx, bookData.ID.String())
	if err != nil {
		s.Logger.Error("service::DeleteBook - failed to delete book: ", err)
		return err
	}

	return nil
}

func (s *BookService) SearchBooks(ctx context.Context, req *dto.SearchBookRequest) (*dto.GetListBookResponse, error) {
	pageSize := req.Limit
	pageIndex := (req.Page - 1) * req.Limit

	booksData, err := s.BookRepo.SearchBooks(ctx, &req.Title, &req.CategoryID, &req.AuthorID, pageSize, pageIndex)
	if err != nil {
		s.Logger.Error("service::SearchBooks - failed to search books: ", err)
		return nil, err
	}

	books := make([]dto.Book, 0)
	for _, book := range booksData {
		books = append(books, dto.Book{
			ID:            book.ID.String(),
			Title:         book.Title,
			Description:   book.Description,
			Isbn:          *book.Isbn,
			PublishedDate: book.PublishedDate.Format(constants.DateTimeFormat),
		})
	}

	pagination := dto.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	response := &dto.GetListBookResponse{
		BookList:   books,
		Pagination: pagination,
	}

	return response, nil
}

func (S *BookService) GetRecommendations(ctx context.Context, userID string, limit, offset int) (*dto.GetListRecommendationsResponse, error) {
	pageSize := limit
	pageIndex := (offset - 1) * limit

	booksData, err := S.BookRepo.GetRecommendations(ctx, userID, pageSize, pageIndex)
	if err != nil {
		S.Logger.Error("service::GetRecommendations - failed to get recommendations: ", err)
		return nil, err
	}

	recommendations := make([]dto.Recommendations, 0)
	for _, book := range booksData {
		recommendations = append(recommendations, dto.Recommendations{
			ID:            book.ID.String(),
			Title:         book.Title,
			AuthorID:      book.AuthorID.String(),
			CategoryID:    book.CategoryID.String(),
			Description:   book.Description,
			PublishedDate: book.PublishedDate.Format(constants.DateTimeFormat),
		})
	}

	pagination := dto.Pagination{
		Page:  offset,
		Limit: limit,
	}

	response := &dto.GetListRecommendationsResponse{
		RecommendationList: recommendations,
		Pagination:         pagination,
	}

	return response, nil
}
