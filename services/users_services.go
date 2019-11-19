package services

import (
	"github.com/sauravgsh16/bookstore_users-api/domain/users"
	"github.com/sauravgsh16/bookstore_users-api/utils/errors"
)

// CreateUser creates a new user in the database
func CreateUser(u users.User) (*users.User, *errors.RestErr) {
	if err := u.Validate(); err != nil {
		return nil, errors.NewBadRequestError("invalid user data")
	}

	return nil, nil
}
