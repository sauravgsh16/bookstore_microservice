package services

import (
	"fmt"

	"github.com/graphql-go/graphql"
	// "github.com/sauravgsh16/bookstore_users-api/utils/errors"
)

// GraphQLResolvers interface
type GraphQLResolvers interface {
	UserResolverFunc(p graphql.ResolveParams) (interface{}, error)
	UsersResolverFunc(p graphql.ResolveParams) (interface{}, error)
}

// Resolver struct
type Resolver struct{}

// UserResolverFunc defines resolver for get user
func (r *Resolver) UserResolverFunc(p graphql.ResolveParams) (interface{}, error) {
	user, err := UserServ.GetUser(p.Args["id"].(int))
	if err != nil {
		return nil, fmt.Errorf(err.Error)
	}
	return user, nil
}

// UsersResolverFunc defines resolver tp get all users with status
func (r *Resolver) UsersResolverFunc(p graphql.ResolveParams) (interface{}, error) {
	users, err := UserServ.SearchUser(p.Args["status"].(string))
	if err != nil {
		return nil, fmt.Errorf(err.Error)
	}

	return users, nil
}
