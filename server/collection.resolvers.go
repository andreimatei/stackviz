package server

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"context"
	"stacksviz/ent"
	"stacksviz/ent/collection"
)

// CreateCollection is the resolver for the createCollection field.
func (r *mutationResolver) CreateCollection(ctx context.Context, input *ent.CreateCollectionInput) (*ent.Collection, error) {
	return r.client.Collection.Create().SetInput(*input).Save(ctx)
}

// CollectionByID is the resolver for the collectionByID field.
func (r *queryResolver) CollectionByID(ctx context.Context, id int) (*ent.Collection, error) {
	return r.client.Debug().Collection.Query().Where(collection.ID(id)).Only(ctx)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }