package graphql

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	schema "github.com/sauravgsh16/bookstore_users-api/domain/graphql-schema"
	"github.com/sauravgsh16/bookstore_users-api/services"
)

// Handler graphql handler
func Handler() gin.HandlerFunc {
	schema.InitQL(&services.Resolver{})

	h := handler.New(&handler.Config{
		Schema:   &schema.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
