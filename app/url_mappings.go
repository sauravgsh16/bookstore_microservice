package app

import (
	"github.com/sauravgsh16/bookstore_users-api/controllers/graphql"
	"github.com/sauravgsh16/bookstore_users-api/controllers/ping"
	"github.com/sauravgsh16/bookstore_users-api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.GET("/users/:user_id", users.Get)
	router.POST("/users", users.Create)
	router.PUT("/users/:user_id", users.Update)
	router.PATCH("/users/:user_id", users.Update)
	router.DELETE("/users/:user_id", users.Delete)
	router.GET("/internal/users/search", users.Search)

	// GraphQL
	router.GET("/graphql", graphql.Handler())
	router.POST("/graphql", graphql.Handler())
}
