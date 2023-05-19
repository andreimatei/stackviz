package server

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"stacksviz/ent"
	"stacksviz/ent/collection"
	"time"
)

// CreateCollection is the resolver for the createCollection field.
func (r *mutationResolver) CreateCollection(ctx context.Context, input *ent.CreateCollectionInput) (*ent.Collection, error) {
	return r.dbClient.Collection.Create().SetInput(*input).Save(ctx)
}

// CollectCollection is the resolver for the collectCollection field.
func (r *mutationResolver) CollectCollection(ctx context.Context) (*ent.Collection, error) {
	i := 0
	psIDs := make([]int, 0, len(r.conf.Targets))

	// Read the name of the first (and only) service.
	if len(r.conf.Targets) != 1 {
		return nil, fmt.Errorf("expected exactly one service")
	}
	var svcName string
	for serviceName, _ := range r.conf.Targets {
		svcName = serviceName
	}

	for processName, url := range r.conf.Targets[svcName] {
		i++
		log.Printf("collecting snapshot from process %d: %s-%s - %s", i, svcName, processName, url)
		resp, err := http.Get(url)
		// TODO(andrei): try the other processes instead of bailing out
		if err != nil {
			return nil, err
		}
		body := resp.Body
		snap, err := io.ReadAll(body)
		if err != nil {
			return nil, err
		}
		ps, err := r.dbClient.ProcessSnapshot.Create().SetInput(ent.CreateProcessSnapshotInput{
			ProcessID: processName,
			Snapshot:  string(snap),
		}).Save(ctx)
		if err != nil {
			return nil, err
		}
		psIDs = append(psIDs, ps.ID)
		// TODO(andrei): close this even in the error cases above.
		_ = body.Close()
	}

	const timeFormat = "Monday, 02-Jan-06 15:04:05 MST"
	return r.dbClient.Collection.Create().SetInput(ent.CreateCollectionInput{
		Name:               fmt.Sprintf("%s - %s", svcName, time.Now().Format(timeFormat)),
		ProcessSnapshotIDs: psIDs,
	}).Save(ctx)
}

// CollectionByID is the resolver for the collectionByID field.
func (r *queryResolver) CollectionByID(ctx context.Context, id int) (*ent.Collection, error) {
	return r.dbClient.Debug().Collection.Query().Where(collection.ID(id)).Only(ctx)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
