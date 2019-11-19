// DAO - domain access object: Provides the means to access the persistance layers
// No other layer in the application is responsible for accessing the db

package users

import (
	"fmt"

	"github.com/sauravgsh16/bookstore_users-api/utils/dates"
	"github.com/sauravgsh16/bookstore_users-api/utils/errors"
)

var (
	userDB = make(map[int64]*User)
)

// Save the user to the db
func (u *User) Save() *errors.RestErr {
	result := userDB[u.ID]
	if result != nil {
		if result.Email == u.Email {
			return errors.NewBadRequestError(fmt.Sprintf("email %s is already registered", u.Email))
		}
		return errors.NewBadRequestError("User already exists")
	}

	u.DateCreated = dates.GetNowString()

	userDB[u.ID] = u
	return nil
}

// Get populates the user pointer or returns error if not found
func (u *User) Get(userID int64) *errors.RestErr {
	result, ok := userDB[userID]
	if !ok {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", userID))
	}

	u.ID = result.ID
	u.FirstName = result.FirstName
	u.LastName = result.LastName
	u.Email = result.Email
	u.DateCreated = result.DateCreated

	return nil
}
