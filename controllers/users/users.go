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

func getUserID(idStr string) (int64, *errors.RestErr) {
	uid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.NewBadRequestError("user id should be a number")
	}
	return uid, nil
}

// Get returns a user
func Get(c *gin.Context) {
	userID, err := getUserID(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	user, getErr := services.GetUser(userID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Create creates a new user
func Create(c *gin.Context) {
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

// Update updates a user
func Update(c *gin.Context) {
	userID, err := getUserID(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	var newUser users.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		bdErr := errors.NewBadRequestError(fmt.Sprintf("invalid json body %s", err.Error()))
		c.JSON(bdErr.Status, bdErr)
		return
	}

	newUser.ID = userID
	isPartial := c.Request.Method == http.MethodPatch

	result, updateErr := services.UpdateUser(newUser, isPartial)
	if err != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

// Delete a user from db
func Delete(c *gin.Context) {
	userID, err := getUserID(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	if err := services.DeleteUser(userID); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// Search searches all users
func Search(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implement me!!\n")
}
