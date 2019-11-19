package app

import (
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

// StartApp starts the user service application
func StartApp() {
	mapUrls()
	router.Run(":8080")
}
