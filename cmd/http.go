package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/hilmiikhsan/library-book-service/helpers"
	"github.com/sirupsen/logrus"
)

func ServeHTTP() {
	// dependency := dependencyInject()

	router := gin.Default()

	err := router.Run(":" + helpers.GetEnv("PORT", ""))
	if err != nil {
		helpers.Logger.Fatal("failed to run http server: ", err)
	}
}

type Dependency struct {
	Logger *logrus.Logger
}

func dependencyInject() Dependency {
	helpers.SetupLogger()

	return Dependency{
		Logger: helpers.Logger,
	}
}
