package book

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/helpers"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/interfaces"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/hilmiikhsan/library-book-service/internal/validator"
)

type BookHandler struct {
	BookService interfaces.IBookService
	Validator   *validator.Validator
}

func (api *BookHandler) CreateBook(ctx *gin.Context) {
	var (
		req = new(dto.CreateBookRequest)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.Logger.Error("handler::CreateBook - Failed to bind request : ", err)
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrFailedBadRequest))
		return
	}

	if err := api.Validator.Validate(req); err != nil {
		helpers.Logger.Error("handler::CreateBook - Failed to validate request : ", err)
		code, errs := helpers.Errors(err, req)
		ctx.JSON(code, helpers.Error(errs))
		return
	}

	err := api.BookService.CreateBook(ctx.Request.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrAuthorNotFound) {
			helpers.Logger.Error("handler::CreateBook - author not found")
			ctx.JSON(http.StatusNotFound, helpers.Error(constants.ErrAuthorNotFound))
			return
		}

		if strings.Contains(err.Error(), constants.ErrCategoryNotFound) {
			helpers.Logger.Error("handler::CreateBook - category not found")
			ctx.JSON(http.StatusNotFound, helpers.Error(constants.ErrCategoryNotFound))
			return
		}

		if strings.Contains(err.Error(), constants.ErrIsbnAlreadyExist) {
			helpers.Logger.Error("handler::CreateBook - isbn already exist")
			ctx.JSON(http.StatusConflict, helpers.Error(constants.ErrIsbnAlreadyExist))
			return
		}

		helpers.Logger.Error("handler::CreateBook - Failed to create Book : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, helpers.Success(nil, ""))
}

func (api *BookHandler) GetDetailBook(ctx *gin.Context) {
	var (
		id = ctx.Param("id")
	)

	if id == "" {
		helpers.Logger.Error("handler::GetDetailBook - Missing required parameter: id")
		ctx.JSON(http.StatusBadRequest, helpers.Error("missing required parameter: id"))
		return
	}

	if !helpers.IsValidUUID(id) {
		helpers.Logger.Error("handler::GetDetailBook - Invalid UUID format for parameter: id")
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrParamIdIsRequired))
		return
	}

	res, err := api.BookService.GetDetailBook(ctx.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrBookNotFound) {
			helpers.Logger.Error("handler::GetDetailBook - Book not found")
			ctx.JSON(http.StatusNotFound, helpers.Error(constants.ErrBookNotFound))
			return
		}

		helpers.Logger.Error("handler::GetDetailBook - Failed to get Book detail : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.Success(res, ""))
}

func (api *BookHandler) GetListBook(ctx *gin.Context) {
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

	res, err := api.BookService.GetListBook(ctx.Request.Context(), pageSize, pageIndex)
	if err != nil {
		helpers.Logger.Error("handler::GetListBook - Failed to get list Book : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.Success(res, ""))
}

func (api *BookHandler) UpdateBook(ctx *gin.Context) {
	var (
		req = new(dto.UpdateBookRequest)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.Logger.Error("handler::UpdateBook - Failed to bind request : ", err)
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrFailedBadRequest))
		return
	}

	if err := api.Validator.Validate(req); err != nil {
		helpers.Logger.Error("handler::UpdateBook - Failed to validate request : ", err)
		code, errs := helpers.Errors(err, req)
		ctx.JSON(code, helpers.Error(errs))
		return
	}

	if !helpers.IsValidUUID(req.ID) {
		helpers.Logger.Error("handler::UpdateStock - Invalid UUID format for parameter: id")
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrIdIsNotValidUUID))
		return
	}

	err := api.BookService.UpdateBook(ctx.Request.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrInvalidFormatDate) {
			helpers.Logger.Error("handler::UpdateBook - Invalid format date")
			ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrInvalidFormatDate))
			return
		}

		if strings.Contains(err.Error(), constants.ErrBookNotFound) {
			helpers.Logger.Error("handler::UpdateBook - book not found")
			ctx.JSON(http.StatusNotFound, helpers.Error(constants.ErrBookNotFound))
			return
		}

		if strings.Contains(err.Error(), constants.ErrAuthorNotFound) {
			helpers.Logger.Error("handler::UpdateBook - author not found")
			ctx.JSON(http.StatusNotFound, helpers.Error(constants.ErrAuthorNotFound))
			return
		}

		if strings.Contains(err.Error(), constants.ErrCategoryNotFound) {
			helpers.Logger.Error("handler::UpdateBook - category not found")
			ctx.JSON(http.StatusNotFound, helpers.Error(constants.ErrCategoryNotFound))
			return
		}

		if strings.Contains(err.Error(), constants.ErrIsbnAlreadyExist) {
			helpers.Logger.Error("handler::UpdateBook - isbn already exist")
			ctx.JSON(http.StatusConflict, helpers.Error(constants.ErrIsbnAlreadyExist))
			return
		}

		helpers.Logger.Error("handler::UpdateBook - Failed to update Book : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.Success(nil, ""))
}

func (api *BookHandler) DeleteBook(ctx *gin.Context) {
	var (
		id = ctx.Param("id")
	)

	if id == "" {
		helpers.Logger.Error("handler::DeleteBook - Missing required parameter: id")
		ctx.JSON(http.StatusBadRequest, helpers.Error("missing required parameter: id"))
		return
	}

	if !helpers.IsValidUUID(id) {
		helpers.Logger.Error("handler::DeleteBook - Invalid UUID format for parameter: id")
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrIdIsNotValidUUID))
		return
	}

	err := api.BookService.DeleteBook(ctx.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrBookNotFound) {
			helpers.Logger.Error("handler::DeleteBook - book not found")
			ctx.JSON(http.StatusNotFound, helpers.Error(constants.ErrBookNotFound))
			return
		}

		helpers.Logger.Error("handler::DeleteBook - Failed to delete book : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.Success(nil, ""))
}

func (api *BookHandler) SearchBooks(ctx *gin.Context) {
	var (
		req = new(dto.SearchBookRequest)
	)

	if err := ctx.ShouldBindJSON(req); err != nil {
		helpers.Logger.Error("handler::SearchBooks - Failed to bind request : ", err)
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrFailedBadRequest))
		return
	}

	if err := api.Validator.Validate(req); err != nil {
		helpers.Logger.Error("handler::SearchBooks - Failed to validate request : ", err)
		code, errs := helpers.Errors(err, req)
		ctx.JSON(code, helpers.Error(errs))
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	res, err := api.BookService.SearchBooks(ctx.Request.Context(), req)
	if err != nil {
		helpers.Logger.Error("handler::SearchBooks - Failed to search books : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.Success(res, ""))
}

func (api *BookHandler) GetRecommendations(ctx *gin.Context) {
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

	token, ok := ctx.Get(constants.TokenTypeAccess)
	if !ok {
		helpers.Logger.Error("handler::GetRecommendations - Failed to get token")
		ctx.JSON(http.StatusUnauthorized, helpers.Error("Failed to get token"))
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		helpers.Logger.Error("handler::GetRecommendations - Failed to parse token")
		ctx.JSON(http.StatusUnauthorized, helpers.Error("Failed to parse token"))
		return
	}

	res, err := api.BookService.GetRecommendations(ctx.Request.Context(), tokenData.UserID, pageSize, pageIndex)
	if err != nil {
		helpers.Logger.Error("handler::GetRecommendations - Failed to get recommendations : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.Success(res, ""))
}
