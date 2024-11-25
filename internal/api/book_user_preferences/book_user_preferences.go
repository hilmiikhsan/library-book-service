package book_user_preferences

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/helpers"
	"github.com/hilmiikhsan/library-book-service/internal/dto"
	"github.com/hilmiikhsan/library-book-service/internal/interfaces"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/hilmiikhsan/library-book-service/internal/validator"
)

type BookUserPreferencesHandler struct {
	BookUserPreferencesService interfaces.IBookUserPreferencesService
	Validator                  *validator.Validator
}

func (api *BookUserPreferencesHandler) CreateBookUserPreferences(ctx *gin.Context) {
	var (
		req = new(dto.CreateBookUserPreferencesRequest)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.Logger.Error("handler::CreateBookUserPreferences - Failed to bind request : ", err)
		ctx.JSON(http.StatusBadRequest, helpers.Error(constants.ErrFailedBadRequest))
		return
	}

	if err := api.Validator.Validate(req); err != nil {
		helpers.Logger.Error("handler::CreateBookUserPreferences - Failed to validate request : ", err)
		code, errs := helpers.Errors(err, req)
		ctx.JSON(code, helpers.Error(errs))
		return
	}

	token, ok := ctx.Get(constants.TokenTypeAccess)
	if !ok {
		helpers.Logger.Error("handler::CreateBookUserPreferences - Failed to get token")
		ctx.JSON(http.StatusUnauthorized, helpers.Error("Failed to get token"))
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		helpers.Logger.Error("handler::CreateBookUserPreferences - Failed to parse token")
		ctx.JSON(http.StatusUnauthorized, helpers.Error("Failed to parse token"))
		return
	}

	req.UserID = tokenData.UserID

	err := api.BookUserPreferencesService.CreateBookUserPreferences(ctx.Request.Context(), req)
	if err != nil {
		helpers.Logger.Error("handler::CreateBookUserPreferences - Failed to create BookUserPreferences : ", err)
		ctx.JSON(http.StatusInternalServerError, helpers.Error(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, helpers.Success(nil, ""))
}
