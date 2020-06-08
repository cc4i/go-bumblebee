package graph

import (
	"context"
	"github.com/gin-gonic/gin"
	"gql/graph/generated"
	"net/http"
	"os"

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
var prefix = "/gql"

func headers(c *gin.Context) {
	ver := os.Getenv("OVERRIDE_VERSION")
	if ver == "" { ver="v1" }
	c.Header("gql_server","air")
	c.Header("gql_version", ver)
}

func Router(ctx context.Context) *gin.Engine {

	r := gin.Default()

	r.POST(prefix+"/query", graphqlHandler())

	r.GET(prefix+"/", playgroundHandler())

	r.GET("/ping", func(c *gin.Context) {
		headers(c)
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
		headers(c)
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", prefix+"/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
