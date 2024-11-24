package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/hilmiikhsan/library-book-service/helpers"
	healthCheckAPI "github.com/hilmiikhsan/library-book-service/internal/api/health_check"
	"github.com/hilmiikhsan/library-book-service/internal/interfaces"
	healthCheckServices "github.com/hilmiikhsan/library-book-service/internal/services/health_check"
	"github.com/sirupsen/logrus"
)

func ServeHTTP() {
	dependency := dependencyInject()

	router := gin.Default()

	router.GET("/health", dependency.HealthcheckAPI.HealthcheckHandlerHTTP)

	err := router.Run(":" + helpers.GetEnv("PORT", ""))
	if err != nil {
		helpers.Logger.Fatal("failed to run http server: ", err)
	}
}

type Dependency struct {
	Logger *logrus.Logger

	HealthcheckAPI interfaces.IHealthcheckHandler
}

func dependencyInject() Dependency {
	helpers.SetupLogger()

	healthcheckSvc := &healthCheckServices.Healthcheck{}
	healthcheckAPI := &healthCheckAPI.Healthcheck{
		HealthcheckServices: healthcheckSvc,
	}

	return Dependency{
		Logger:         helpers.Logger,
		HealthcheckAPI: healthcheckAPI,
	}
}
