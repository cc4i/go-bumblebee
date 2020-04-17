package main

import (
	"gql/graph"
	"log"
)

func main() {
	log.Fatal(graph.Router().Run("0.0.0.0:9030"))
}
