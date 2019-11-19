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

	if err := u.Save(); err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUser returns user if present
func GetUser(userID int64) (*users.User, *errors.RestErr) {
	user := &users.User{}
	if err := user.Get(userID); err != nil {
		return nil, err
	}
	return user, nil
}
