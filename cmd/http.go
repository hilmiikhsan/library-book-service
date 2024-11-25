package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/hilmiikhsan/library-book-service/external"
	"github.com/hilmiikhsan/library-book-service/helpers"
	bookAPI "github.com/hilmiikhsan/library-book-service/internal/api/book"
	bookBorrowedAPI "github.com/hilmiikhsan/library-book-service/internal/api/book_borrowed"
	bookStockAPI "github.com/hilmiikhsan/library-book-service/internal/api/book_stock"
	bookUserPreferencesAPI "github.com/hilmiikhsan/library-book-service/internal/api/book_user_preferences"
	healthCheckAPI "github.com/hilmiikhsan/library-book-service/internal/api/health_check"
	"github.com/hilmiikhsan/library-book-service/internal/interfaces"
	bookRepository "github.com/hilmiikhsan/library-book-service/internal/repository/book"
	bookBorrowedRepository "github.com/hilmiikhsan/library-book-service/internal/repository/book_borrowed"
	bookStockRepository "github.com/hilmiikhsan/library-book-service/internal/repository/book_stock"
	bookUserPreferencesRepository "github.com/hilmiikhsan/library-book-service/internal/repository/book_user_preferences"
	bookServices "github.com/hilmiikhsan/library-book-service/internal/services/book"
	bookBorrowedServices "github.com/hilmiikhsan/library-book-service/internal/services/book_borrowed"
	bookStockServices "github.com/hilmiikhsan/library-book-service/internal/services/book_stock"
	bookUserPreferencesServices "github.com/hilmiikhsan/library-book-service/internal/services/book_user_preferences"
	healthCheckServices "github.com/hilmiikhsan/library-book-service/internal/services/health_check"
	"github.com/hilmiikhsan/library-book-service/internal/validator"
	"github.com/sirupsen/logrus"
)

func ServeHTTP() {
	dependency := dependencyInject()

	router := gin.Default()

	router.GET("/health", dependency.HealthcheckAPI.HealthcheckHandlerHTTP)

	bookV1 := router.Group("/book/v1")
	bookV1.POST("/create", dependency.MiddlewareValidateAdminToken, dependency.BookAPI.CreateBook)
	bookV1.GET("/:id", dependency.MiddlewareValidateToken, dependency.BookAPI.GetDetailBook)
	bookV1.GET("/", dependency.MiddlewareValidateToken, dependency.BookAPI.GetListBook)
	bookV1.PUT("/update", dependency.MiddlewareValidateAdminToken, dependency.BookAPI.UpdateBook)
	bookV1.DELETE("/:id", dependency.MiddlewareValidateAdminToken, dependency.BookAPI.DeleteBook)
	bookV1.GET("/search", dependency.MiddlewareValidateUserToken, dependency.BookAPI.SearchBooks)
	bookV1.GET("/recommendations", dependency.MiddlewareValidateUserToken, dependency.BookAPI.GetRecommendations)

	bookStockV1 := router.Group("/book-stock/v1")
	bookStockV1.POST("/create", dependency.MiddlewareValidateAdminToken, dependency.BookStockAPI.CreateBookStock)
	bookStockV1.GET("/:id", dependency.MiddlewareValidateAdminToken, dependency.BookStockAPI.GetDetailBookStock)
	bookStockV1.GET("/", dependency.MiddlewareValidateAdminToken, dependency.BookStockAPI.GetListBookStock)
	bookStockV1.PUT("/update", dependency.MiddlewareValidateAdminToken, dependency.BookStockAPI.UpdateBookStock)
	bookStockV1.DELETE("/:id", dependency.MiddlewareValidateAdminToken, dependency.BookStockAPI.DeleteBookStock)

	bookBorrowedV1 := router.Group("/book-borrowed/v1")
	bookBorrowedV1.POST("/borrow", dependency.MiddlewareValidateUserToken, dependency.BookBorrowedAPI.BookBorrowed)
	bookBorrowedV1.POST("/return", dependency.MiddlewareValidateUserToken, dependency.BookBorrowedAPI.BookReturned)

	bookUserPreferencesV1 := router.Group("/book-user-preferences/v1")
	bookUserPreferencesV1.POST("/create", dependency.MiddlewareValidateUserToken, dependency.BookUserPreferencesAPI.CreateBookUserPreferences)

	err := router.Run(":" + helpers.GetEnv("PORT", ""))
	if err != nil {
		helpers.Logger.Fatal("failed to run http server: ", err)
	}
}

type Dependency struct {
	Logger                        *logrus.Logger
	BookRepository                interfaces.IBookRepository
	BookStockRepository           interfaces.IBookStockRepository
	BookBorrowedRepository        interfaces.IBookBorrowedRepository
	BookUserPreferencesRepository interfaces.IBookUserPreferencesRepository

	HealthcheckAPI         interfaces.IHealthcheckHandler
	BookAPI                interfaces.IBookHandler
	BookStockAPI           interfaces.IBookStockHandler
	BookBorrowedAPI        interfaces.IBookBorrowedHandler
	BookUserPreferencesAPI interfaces.IBookUserPreferencesHandler
	External               interfaces.IExternal
}

func dependencyInject() Dependency {
	helpers.SetupLogger()

	healthcheckSvc := &healthCheckServices.Healthcheck{}
	healthcheckAPI := &healthCheckAPI.Healthcheck{
		HealthcheckServices: healthcheckSvc,
	}

	bookRepo := &bookRepository.BookRepository{
		DB:     helpers.DB,
		Logger: helpers.Logger,
		Redis:  helpers.RedisClient,
	}

	bookStockRepo := &bookStockRepository.BookStockRepository{
		DB:     helpers.DB,
		Logger: helpers.Logger,
		Redis:  helpers.RedisClient,
	}

	bookBorrowedRepo := &bookBorrowedRepository.BookBorrowedRepository{
		DB:     helpers.DB,
		Logger: helpers.Logger,
	}

	bookUserPreferencesRepo := &bookUserPreferencesRepository.BookUserPreferencesRepository{
		DB:     helpers.DB,
		Logger: helpers.Logger,
	}

	validator := validator.NewValidator()

	external := &external.External{
		Logger: helpers.Logger,
	}

	bookSvc := &bookServices.BookService{
		BookRepo: bookRepo,
		External: external,
		Logger:   helpers.Logger,
	}
	bookAPI := &bookAPI.BookHandler{
		BookService: bookSvc,
		Validator:   validator,
	}

	bookStockSvc := &bookStockServices.BookStockService{
		BookStockRepo: bookStockRepo,
		BookRepo:      bookRepo,
		Logger:        helpers.Logger,
	}
	bookStockAPI := &bookStockAPI.BookStockHandler{
		BookStockService: bookStockSvc,
		Validator:        validator,
	}

	bookBorrowedSvc := &bookBorrowedServices.BookBorrowedService{
		BookBorrowedRepo: bookBorrowedRepo,
		BookStockRepo:    bookStockRepo,
		Logger:           helpers.Logger,
		DB:               helpers.DB,
	}
	bookBorrowedAPI := &bookBorrowedAPI.BookBorrowedHandler{
		BookBorrowedService: bookBorrowedSvc,
		Validator:           validator,
	}

	bookUserPreferencesSvc := &bookUserPreferencesServices.BookUserPreferencesService{
		BookUserPreferencesRepo: bookUserPreferencesRepo,
		External:                external,
		Logger:                  helpers.Logger,
	}
	bookUserPreferencesAPI := &bookUserPreferencesAPI.BookUserPreferencesHandler{
		BookUserPreferencesService: bookUserPreferencesSvc,
		Validator:                  validator,
	}

	return Dependency{
		Logger:                        helpers.Logger,
		BookRepository:                bookRepo,
		BookStockRepository:           bookStockRepo,
		BookBorrowedRepository:        bookBorrowedRepo,
		BookUserPreferencesRepository: bookUserPreferencesRepo,
		HealthcheckAPI:                healthcheckAPI,
		BookAPI:                       bookAPI,
		BookStockAPI:                  bookStockAPI,
		BookBorrowedAPI:               bookBorrowedAPI,
		BookUserPreferencesAPI:        bookUserPreferencesAPI,
		External:                      external,
	}
}
