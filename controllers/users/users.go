package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sauravgsh16/bookstore_users-api/domain/users"
	"github.com/sauravgsh16/bookstore_users-api/services"
	"github.com/sauravgsh16/bookstore_users-api/utils/errors"
)

// GetUser returns a user
func GetUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		invErr := errors.NewBadRequestError("user id should be a number")
		c.JSON(invErr.Status, invErr)
		return
	}

	user, getErr := services.GetUser(userID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var user users.User

	// ShouldBindJSON - read request body and unmarshals the []bytes to user
	if err := c.ShouldBindJSON(&user); err != nil {
		bdErr := errors.NewBadRequestError(fmt.Sprintf("invalid json body: %s", err.Error()))
		c.JSON(bdErr.Status, bdErr)
		return
	}

	result, err := services.CreateUser(user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// SearchUser searches all users
func SearchUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implement me!!\n")
}
