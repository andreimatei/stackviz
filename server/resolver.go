package server

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"stacksviz/datasource"
	"stacksviz/ent"
	"stacksviz/util"

	"github.com/99designs/gqlgen/graphql"
)

// Resolver is the resolver root.
type Resolver struct {
	stacksFetcher datasource.StacksFetcher
	dbClient      *ent.Client
	conf          util.Config
}

// NewSchema creates a graphql executable schema.
func NewSchema(client *ent.Client, stacksFetcher datasource.StacksFetcher, conf util.Config) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{
			stacksFetcher: stacksFetcher,
			dbClient:      client.Debug(), // client.Debug() to log all the queries
			conf:          conf,
		},
	})
}
