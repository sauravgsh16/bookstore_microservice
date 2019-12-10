package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sauravgsh16/bookstore_users-api/utils/crypto"

	"github.com/gin-gonic/gin"
	"github.com/sauravgsh16/bookstore_users-api/domain/users"
	"github.com/sauravgsh16/bookstore_users-api/services"
	"github.com/sauravgsh16/bookstore_users-api/utils/errors"
)

func getUserID(idStr string) (int, *errors.RestErr) {
	uid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.NewBadRequestError("user id should be a number")
	}
	return int(uid), nil
}

// Get returns a user
func Get(c *gin.Context) {
	userID, err := getUserID(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	user, getErr := services.UserServ.GetUser(userID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	isPublic := c.GetHeader("X-Public") == "true"
	c.JSON(http.StatusOK, user.Marshall(isPublic))
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

	result, err := services.UserServ.CreateUser(user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	isPublic := c.GetHeader("X-Public") == "true"
	c.JSON(http.StatusCreated, result.Marshall(isPublic))
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

	result, updateErr := services.UserServ.UpdateUser(newUser, isPartial)
	if err != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}

	isPublic := c.GetHeader("X-Public") == "true"
	c.JSON(http.StatusOK, result.Marshall(isPublic))
}

// Delete a user from db
func Delete(c *gin.Context) {
	userID, err := getUserID(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	if err := services.UserServ.DeleteUser(userID); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// Search searches all users
func Search(c *gin.Context) {
	status := c.Query("status")

	users, err := services.UserServ.SearchUser(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	isPublic := c.GetHeader("X-Public") == "true"
	c.JSON(http.StatusOK, users.Marshall(isPublic))
}

// LoginUser logs in a user
func LoginUser(c *gin.Context) {
	var req users.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rstErr := errors.NewBadRequestError("invalid request body")
		c.JSON(rstErr.Status, rstErr)
		return
	}

	req.Password = crypto.GetMd5(req.Password)

	user, err := services.UserServ.LoginUser(req)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	isPublic := c.GetHeader("X-Public") == "true"
	c.JSON(http.StatusOK, user.Marshall(isPublic))
}
