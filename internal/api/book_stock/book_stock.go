package book_stock

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/helpers"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/interfaces"
	"github.com/hilmiikhsan/library-book-service/internal/validator"
)

type BookStockHandler struct {
	BookStockService interfaces.IBookStockService
	Validator        *validator.Validator
}

func (api *BookStockHandler) CreateBookStock(ctx *gin.Context) {
	var (
		req = new(dto.CreateBookStockRequest)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.Logger.Error("handler::CreateBookStock - Failed to bind request : ", err)
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrFailedBadRequest))
		return
	}

	if err := api.Validator.Validate(req); err != nil {
		helpers.Logger.Error("handler::CreateBookStock - Failed to validate request : ", err)
		code, errs := helpers.Errors(err, req)
		ctx.JSON(code, helpers.Error(errs))
		return
	}

	err := api.BookStockService.CreateBookStock(ctx.Request.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrBookNotFound) {
			helpers.Logger.Error("handler::CreateBookStock - book not found")
			ctx.JSON(http.StatusNotFound, helpers.Error(constants.ErrBookNotFound))
			return
		}

		if strings.Contains(err.Error(), constants.ErrBookStockAlreadyExist) {
			helpers.Logger.Error("handler::CreateBookStock - BookStock already exist")
			ctx.JSON(http.StatusConflict, helpers.Error(constants.ErrBookStockAlreadyExist))
			return
		}

		helpers.Logger.Error("handler::CreateBookStock - Failed to create BookStock : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, helpers.Success(nil, ""))
}

func (api *BookStockHandler) GetDetailBookStock(ctx *gin.Context) {
	var (
		id = ctx.Param("id")
	)

	if id == "" {
		helpers.Logger.Error("handler::GetDetailBookStock - Missing required parameter: id")
		ctx.JSON(http.StatusBadRequest, helpers.Error("missing required parameter: id"))
		return
	}

	if !helpers.IsValidUUID(id) {
		helpers.Logger.Error("handler::GetDetailBookStock - Invalid UUID format for parameter: id")
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrParamIdIsRequired))
		return
	}

	res, err := api.BookStockService.GetDetailBookStock(ctx.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrBookStockNotFound) {
			helpers.Logger.Error("handler::GetDetailBookStock - BookStock not found")
			ctx.JSON(http.StatusNotFound, helpers.Error(constants.ErrBookStockNotFound))
			return
		}

		helpers.Logger.Error("handler::GetDetailBookStock - Failed to get BookStock detail : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.Success(res, ""))
}

func (api *BookStockHandler) GetListBookStock(ctx *gin.Context) {
	pageIndexStr := ctx.Query("page")
	pageSizeStr := ctx.Query("limit")

	pageIndex, _ := strconv.Atoi(pageIndexStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if pageIndex <= 0 {
		pageIndex = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	res, err := api.BookStockService.GetListBookStock(ctx.Request.Context(), pageSize, pageIndex)
	if err != nil {
		helpers.Logger.Error("handler::GetListBookStock - Failed to get list BookStock : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.Success(res, ""))
}

func (api *BookStockHandler) UpdateBookStock(ctx *gin.Context) {
	var (
		req = new(dto.UpdateBookStockRequest)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.Logger.Error("handler::UpdateBookStock - Failed to bind request : ", err)
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrFailedBadRequest))
		return
	}

	if err := api.Validator.Validate(req); err != nil {
		helpers.Logger.Error("handler::UpdateBookStock - Failed to validate request : ", err)
		code, errs := helpers.Errors(err, req)
		ctx.JSON(code, helpers.Error(errs))
		return
	}

	if !helpers.IsValidUUID(req.ID) {
		helpers.Logger.Error("handler::UpdateStock - Invalid UUID format for parameter: id")
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrIdIsNotValidUUID))
		return
	}

	err := api.BookStockService.UpdateBookStock(ctx.Request.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrBookStockNotFound) {
			helpers.Logger.Error("handler::UpdateBookStock - BookStock not found")
			ctx.JSON(http.StatusNotFound, helpers.Error(constants.ErrBookStockNotFound))
			return
		}

		helpers.Logger.Error("handler::UpdateBookStock - Failed to update BookStock : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.Success(nil, ""))
}

func (api *BookStockHandler) DeleteBookStock(ctx *gin.Context) {
	var (
		id = ctx.Param("id")
	)

	if id == "" {
		helpers.Logger.Error("handler::DeleteBookStock - Missing required parameter: id")
		ctx.JSON(http.StatusBadRequest, helpers.Error("missing required parameter: id"))
		return
	}

	if !helpers.IsValidUUID(id) {
		helpers.Logger.Error("handler::DeleteBookStock - Invalid UUID format for parameter: id")
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrIdIsNotValidUUID))
		return
	}

	err := api.BookStockService.DeleteBookStock(ctx.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrBookStockNotFound) {
			helpers.Logger.Error("handler::DeleteBookStock - BookStock not found")
			ctx.JSON(http.StatusNotFound, helpers.Error(constants.ErrBookStockNotFound))
			return
		}

		helpers.Logger.Error("handler::DeleteBookStock - Failed to delete BookStock : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.Success(nil, ""))
}
