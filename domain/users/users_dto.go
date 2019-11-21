// DTO - domain transfer object - Provides the definitions of the database objects

package users

import (
	"strings"
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
func (u *User) Validate() bool {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if len(u.Email) == 0 {
		return false
	}
	return true
}
