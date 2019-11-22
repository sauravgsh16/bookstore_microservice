package services

import (
	"github.com/sauravgsh16/bookstore_users-api/domain/users"
	"github.com/sauravgsh16/bookstore_users-api/utils/crypto"
	"github.com/sauravgsh16/bookstore_users-api/utils/dates"
	"github.com/sauravgsh16/bookstore_users-api/utils/errors"
)

// CreateUser creates a new user in the database
func CreateUser(u users.User) (*users.User, *errors.RestErr) {
	if valid := u.Validate(); !valid {
		return nil, errors.NewBadRequestError("invalid user data")
	}

	u.DateCreated = dates.GetNowDBString()
	u.Status = users.StatusActive
	u.Password = crypto.GetMd5(u.Password)

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

// UpdateUser updates a user
func UpdateUser(u users.User, isPatch bool) (*users.User, *errors.RestErr) {
	current, err := GetUser(u.ID)
	if err != nil {
		return nil, err
	}

	if isPatch {
		if u.FirstName != "" {
			current.FirstName = u.FirstName
		}

		if u.LastName != "" {
			current.LastName = u.LastName
		}

		if u.Email != "" {
			current.Email = u.Email
		}
	} else {
		current.FirstName = u.FirstName
		current.LastName = u.LastName
		current.Email = u.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

// DeleteUser api
func DeleteUser(uid int64) *errors.RestErr {
	user, err := GetUser(uid)
	if err != nil {
		return err
	}

	if err := user.Delete(); err != nil {
		return err
	}
	return nil
}

// Search returns users matching passed argument
func Search(status string) (users.Users, *errors.RestErr) {
	dao := users.User{}
	return dao.FindByStatus(status)
}
