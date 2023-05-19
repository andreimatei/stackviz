package server

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"stacksviz/ent"
	"stacksviz/util"

	"github.com/99designs/gqlgen/graphql"
)

// Resolver is the resolver root.
type Resolver struct {
	dbClient *ent.Client
	conf     util.Config
}

// NewSchema creates a graphql executable schema.
func NewSchema(client *ent.Client, conf util.Config) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{
			dbClient: client,
			conf:     conf,
		},
	})
}
