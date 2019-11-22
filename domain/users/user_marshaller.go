package users

import (
	"encoding/json"
	"fmt"
)

// Marshaller interface
type Marshaller interface {
	IsMarshalled() bool
}

// PublicUser for external usage
type PublicUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Status    string `json:"status"`
}

// IsMarshalled returns true if it can be marshalled
func (pu PublicUser) IsMarshalled() bool {
	return true
}

// PrivateUser for internal usage
type PrivateUser struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

// IsMarshalled returns true if it can be marshalled
func (pu PrivateUser) IsMarshalled() bool {
	return true
}

// Marshall maps user to either public or private user
func (u User) Marshall(isPublic bool) Marshaller {
	uJSON, err := json.Marshal(u)
	if err != nil {
		fmt.Printf("Error while Marshalling\n")
		return nil
	}
	if isPublic {
		var pubUser PublicUser
		if err := json.Unmarshal(uJSON, &pubUser); err != nil {
			fmt.Printf("Error while unmarshalling: %s\n", err.Error())
			return nil
		}
		return pubUser
	}
	var priUser PrivateUser
	if err := json.Unmarshal(uJSON, &priUser); err != nil {
		fmt.Printf("Error while unmarshalling: %s\n", err.Error())
		return nil
	}
	return priUser
}

// Marshall returns a slice of Marshallers
func (u Users) Marshall(isPublic bool) []Marshaller {
	result := make([]Marshaller, len(u))
	for i, user := range u {
		result[i] = user.Marshall(isPublic)
	}
	return result
}
