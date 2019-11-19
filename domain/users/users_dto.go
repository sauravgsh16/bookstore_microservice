// DTO - domain transfer object - Provides the definitions of the database objects

package users

import (
	"strings"
	"time"

	"github.com/sauravgsh16/bookstore_users-api/utils/errors"
)

// User struct
type User struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
}

// Validate if the user fields are accepted
func (u *User) Validate() *errors.RestErr {
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if len(u.Email) == 0 {
		return nil
	}

	time.Now()
	return nil
}
