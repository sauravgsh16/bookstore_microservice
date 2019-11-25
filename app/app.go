package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sauravgsh16/bookstore_users-api/logger"
)

var (
	router = gin.Default()
)

// StartApp starts the user service application
func StartApp() {
	mapUrls()

	logger.Info("about to start application....")
	if err := router.Run(":8080"); err != nil {
		logger.Error("failed to run gin gonic server, error: ", err)
		panic(err)
	}
}
