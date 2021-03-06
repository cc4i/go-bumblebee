package graph

import (
	"github.com/gin-gonic/gin"
	"gql/graph/generated"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// !!!
// To re-generate graphql module as yaml, for more detail: https://gqlgen.com/getting-started/
//
// $ cd graphqlsvr
// $ go run github.com/99designs/gqlgen generate
//
/// !!!

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/query", graphqlHandler())

	r.GET("/", playgroundHandler())

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
