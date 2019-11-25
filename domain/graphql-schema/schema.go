package schema

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/sauravgsh16/bookstore_users-api/logger"
	"github.com/sauravgsh16/bookstore_users-api/services"
)

// Schema graphql schema
var Schema graphql.Schema

// InitQL schema
func InitQL(r services.GraphQLResolvers) {
	var userType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"first_name": &graphql.Field{
				Type: graphql.String,
			},
			"last_name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"date_created": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	var usersType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Users",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
			},
		},
	})

	fields := graphql.Fields{
		"User": &graphql.Field{
			Type: graphql.Type(userType),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: r.UserResolverFunc,
		},
		"Users": &graphql.Field{
			Type: graphql.NewList(usersType),
			Args: graphql.FieldConfigArgument{
				"status": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.UsersResolverFunc,
		},
	}

	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}

	var err error
	Schema, err = graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create new schema: %s", err.Error()))
	}
	logger.Info("Successfully initialized GraphQL")
}
