package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping responds with Pong
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Pong",
	})
}
