package borrowed_book

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/helpers"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/interfaces"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/hilmiikhsan/library-book-service/internal/validator"
)

type BookBorrowedHandler struct {
	BookBorrowedService interfaces.IBookBorrowedService
	Validator           *validator.Validator
}

func (api *BookBorrowedHandler) BookBorrowed(ctx *gin.Context) {
	var (
		req = new(dto.BookBorrowedRequest)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.Logger.Error("handler::BookBorrowed - Failed to bind request : ", err)
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrFailedBadRequest))
		return
	}

	if err := api.Validator.Validate(req); err != nil {
		helpers.Logger.Error("handler::BookBorrowed - Failed to validate request : ", err)
		code, errs := helpers.Errors(err, req)
		ctx.JSON(code, helpers.Error(errs))
		return
	}

	token, ok := ctx.Get(constants.TokenTypeAccess)
	if !ok {
		helpers.Logger.Error("handler::BookBorrowed - Failed to get token")
		ctx.JSON(http.StatusUnauthorized, helpers.Error("Failed to get token"))
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		helpers.Logger.Error("handler::BookBorrowed - Failed to parse token")
		ctx.JSON(http.StatusUnauthorized, helpers.Error("Failed to parse token"))
		return
	}

	err := api.BookBorrowedService.BookBorrowed(ctx.Request.Context(), req, tokenData.UserID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrInvalidFormatDate) {
			helpers.Logger.Error("handler::BookBorrowed - Invalid format date : ", err)
			ctx.JSON(http.StatusBadRequest, helpers.Error(err.Error()))
			return
		}

		if strings.Contains(err.Error(), constants.ErrBookStockNotFound) {
			helpers.Logger.Error("handler::BookBorrowed - Book stock not found : ", err)
			ctx.JSON(http.StatusNotFound, helpers.Error(err.Error()))
			return
		}

		if strings.Contains(err.Error(), constants.ErrBookAlreadyBorrowed) {
			helpers.Logger.Error("handler::BookBorrowed - Book already borrowed : ", err)
			ctx.JSON(http.StatusConflict, helpers.Error(err.Error()))
			return
		}

		if strings.Contains(err.Error(), constants.ErrInsufficientStock) {
			helpers.Logger.Error("handler::BookBorrowed - Insufficient stock : ", err)
			ctx.JSON(http.StatusConflict, helpers.Error(err.Error()))
			return
		}

		helpers.Logger.Error("handler::BookBorrowed - Failed to borrow book : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, helpers.Success(nil, ""))
}

func (api *BookBorrowedHandler) BookReturned(ctx *gin.Context) {
	var (
		req = new(dto.BookReturnedRequest)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.Logger.Error("handler::BookReturned - Failed to bind request : ", err)
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrFailedBadRequest))
		return
	}

	if err := api.Validator.Validate(req); err != nil {
		helpers.Logger.Error("handler::BookReturned - Failed to validate request : ", err)
		code, errs := helpers.Errors(err, req)
		ctx.JSON(code, helpers.Error(errs))
		return
	}

	token, ok := ctx.Get(constants.TokenTypeAccess)
	if !ok {
		helpers.Logger.Error("handler::BookReturned - Failed to get token")
		ctx.JSON(http.StatusUnauthorized, helpers.Error("Failed to get token"))
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		helpers.Logger.Error("handler::BookReturned - Failed to parse token")
		ctx.JSON(http.StatusUnauthorized, helpers.Error("Failed to parse token"))
		return
	}

	err := api.BookBorrowedService.BookReturned(ctx.Request.Context(), req, tokenData.UserID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrInvalidFormatDate) {
			helpers.Logger.Error("handler::BookBorrowed - Invalid format date : ", err)
			ctx.JSON(http.StatusBadRequest, helpers.Error(err.Error()))
			return
		}

		if strings.Contains(err.Error(), constants.ErrBookStockNotFound) {
			helpers.Logger.Error("handler::BookReturned - Book borrowed not found : ", err)
			ctx.JSON(http.StatusNotFound, helpers.Error(err.Error()))
			return
		}

		if strings.Contains(err.Error(), constants.ErrBookAlreadyReturned) {
			helpers.Logger.Error("handler::BookReturned - Book already returned : ", err)
			ctx.JSON(http.StatusConflict, helpers.Error(err.Error()))
			return
		}

		helpers.Logger.Error("handler::BookReturned - Failed to return book : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.Success(nil, ""))
}
