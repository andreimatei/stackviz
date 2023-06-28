package server

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"stacksviz/ent"
	"stacksviz/ent/collection"
	"stacksviz/ent/framespec"
	"stacksviz/stacks"
	"time"

	"golang.org/x/sync/errgroup"
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
		// TODO(andrei): deal with multiple services
		svcName = serviceName
	}

	var g errgroup.Group
	spec := r.getOrCreateCollectSpec(ctx)
	for processName, url := range r.conf.Targets[svcName] {
		i++
		processName := processName
		url := url
		g.Go(func() error {
			log.Printf("collecting snapshot from process %d: %s-%s - %s", i, svcName, processName, url)
			// !!! snap, err := r.getSnapshotFromPprof(url)
			snap, err := r.getSnapshotFromDelveAgent(ctx, url, spec)
			if err != nil {
				return err
			}
			log.Printf("!!! delve agent returned")
			var framesOfInterest string
			if len(snap.Frames_of_interest) > 0 {
				b, err := json.Marshal(snap.Frames_of_interest)
				if err != nil {
					return err
				}
				framesOfInterest = string(b)
			}

			log.Printf("!!! creating snapshot with frames of interest: %s", framesOfInterest)
			input := ent.CreateProcessSnapshotInput{
				ProcessID: processName,
				Snapshot:  snapToString(snap),
			}
			if framesOfInterest != "" {
				input.FramesOfInterest = &framesOfInterest
			}
			ps, err := r.dbClient.ProcessSnapshot.Create().SetInput(input).Save(ctx)
			if err != nil {
				return err
			}
			psIDs = append(psIDs, ps.ID)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	const timeFormat = "Monday, 02-Jan-06 15:04:05 MST"
	return r.dbClient.Collection.Create().SetInput(ent.CreateCollectionInput{
		Name:               fmt.Sprintf("%s - %s", svcName, time.Now().Format(timeFormat)),
		ProcessSnapshotIDs: psIDs,
	}).Save(ctx)
}

// AddExprToCollectSpec is the resolver for the addExprToCollectSpec mutation.
func (r *mutationResolver) AddExprToCollectSpec(ctx context.Context, frame string, expr string) (*ent.CollectSpec, error) {
	// TODO(andrei): use transactions
	ci := r.getOrCreateCollectSpec(ctx)
	// Create of update the frame info.
	fi, err := ci.QueryFrames().Where(framespec.Frame(frame)).Only(ctx)
	nfe := &ent.NotFoundError{}
	if errors.As(err, &nfe) {
		fi = r.dbClient.FrameSpec.Create().SetFrame(frame).SetExprs([]string{expr}).SaveX(ctx)
		ci = ci.Update().AddFrames(fi).SaveX(ctx)
		return ci, nil
	}
	if err != nil {
		return nil, err
	}
	for _, e := range fi.Exprs {
		if e == expr {
			return ci, nil
		}
	}
	fi.Update().SetExprs(append(fi.Exprs, expr)).SaveX(ctx)
	return ci, nil
}

// RemoveExprFromCollectSpec is the resolver for the removeExprFromCollectSpec field.
func (r *mutationResolver) RemoveExprFromCollectSpec(ctx context.Context, expr string, frame string) (*ent.CollectSpec, error) {
	ci := r.getOrCreateCollectSpec(ctx)
	fi, err := ci.QueryFrames().Where(framespec.Frame(frame)).Only(ctx)
	nfe := &ent.NotFoundError{}
	if errors.As(err, &nfe) {
		return ci, nil
	}
	foundIndex := -1
	for i, e := range fi.Exprs {
		if e == expr {
			foundIndex = i
			break
		}
	}
	if foundIndex != -1 {
		fi.Update().SetExprs(append(fi.Exprs[:foundIndex], fi.Exprs[foundIndex+1:]...)).SaveX(ctx)
	}
	return ci, nil
}

// CollectionByID is the resolver for the collectionByID field.
func (r *queryResolver) CollectionByID(ctx context.Context, id int) (*ent.Collection, error) {
	log.Printf("!!! querying collection by ID: %d", id)
	return r.dbClient.Debug().Collection.Query().Where(collection.ID(id)).Only(ctx)
}

// Goroutines is the resolver for the goroutines field.
func (r *queryResolver) Goroutines(ctx context.Context, colID int, snapID int, gID *int) (*SnapshotInfo, error) {
	snap, err := r.stacksFetcher.Fetch(ctx, colID, snapID)
	if err != nil {
		return nil, err
	}
	gMap := make(map[int]*GoroutineInfo, len(snap.Snapshot.Goroutines))
	for _, g := range snap.Snapshot.Goroutines {
		frames := make([]*FrameInfo, len(g.Stack.Calls))
		for j, c := range g.Stack.Calls {
			frames[j] = &FrameInfo{
				Func: c.Func.Complete,
				File: c.RemoteSrcPath,
				Line: c.Line,
			}
		}

		// Render all the frames of interest for the goroutine, across all the frames.
		// !!! preprocess an index from variable value to list of links.
		var vs []*CollectedVar
		for _, frames := range snap.FramesOfInterest[g.ID] {
			for _, v := range frames.Vars {
				links := make([]*Link, 0, len(v.Links))
				for _, l := range v.Links {
					if l.GoroutineID == g.ID {
						// Don't link to ourselves.
						continue
					}
					links = append(links,
						&Link{
							SnapshotID:  l.SnapshotID,
							GoroutineID: l.GoroutineID,
							FrameIdx:    l.FrameIdx,
						})
				}
				vs = append(vs, &CollectedVar{
					Value: v.Val,
					Links: links,
				})
			}
		}

		gMap[g.ID] = &GoroutineInfo{
			ID:     g.ID,
			Frames: frames,
			Vars:   vs,
		}
	}

	if gID != nil {
		log.Printf("!!! filtering for goroutine: %d", *gID)
		if gi := gMap[*gID]; gi != nil {
			return &SnapshotInfo{
				Raw:        []*GoroutineInfo{gi},
				Aggregated: nil,
			}, nil
		}
		return nil, nil
	} else {
		log.Printf("!!! not filtering for a specific goroutine")
	}

	allGs := make([]*GoroutineInfo, 0, len(gMap))
	for _, gi := range gMap {
		allGs = append(allGs, gi)
	}

	groups := make([]*GoroutinesGroup, len(snap.Agg.Buckets))
	for i, b := range snap.Agg.Buckets {
		frames := make([]*FrameInfo, len(b.Stack.Calls))
		for j, c := range b.Stack.Calls {
			frames[j] = &FrameInfo{
				Func: c.Func.Complete,
				File: c.RemoteSrcPath,
				Line: c.Line,
			}
		}
		groups[i] = &GoroutinesGroup{
			IDs:    b.IDs,
			Frames: frames,
		}

		for _, gID := range b.IDs {
			if gi := gMap[gID]; gi != nil {
				groups[i].Vars = append(groups[i].Vars, gi.Vars...)
			}
		}
	}

	si := &SnapshotInfo{
		Raw:        allGs,
		Aggregated: groups,
	}
	return si, nil
}

// AvailableVars is the resolver for the availableVars field.
func (r *queryResolver) AvailableVars(ctx context.Context, funcArg string, pcOff int) (*VarsAndTypes, error) {
	var svcName string
	for serviceName, _ := range r.conf.Targets {
		// TODO(andrei): deal with multiple services
		svcName = serviceName
	}

	var agentURL string
	for _, url := range r.conf.Targets[svcName] {
		agentURL = url
		break
	}

	log.Printf("!!! calling AvailableVars on Delve agent")
	vars, types, err := r.getAvailableVarsFromDelveAgent(agentURL, funcArg, int64(pcOff))
	if err != nil {
		return nil, err
	}
	resVars := make([]*VarInfo, len(vars))
	for i, v := range vars {
		resVars[i] = &VarInfo{
			Name:             v.Name,
			Type:             v.Type,
			FormalParameter:  v.FormalParameter,
			LoclistAvailable: v.LoclistAvailable,
		}
	}
	resTypes := make([]*TypeInfo, len(types))
	for i, t := range types {
		resTypes[i] = &TypeInfo{
			Name:            t.Name,
			FieldsNotLoaded: t.FieldsNotLoaded,
		}
		resTypes[i].Fields = make([]*FieldInfo, len(t.Fields))
		for j, f := range t.Fields {
			resTypes[i].Fields[j] = &FieldInfo{
				Name:     f.Name,
				Type:     f.TypeName,
				Embedded: f.Embedded,
			}
		}
	}

	res := &VarsAndTypes{
		Vars:  resVars,
		Types: resTypes,
	}
	return res, nil
}

// FrameInfo is the resolver for the frameInfo field.
func (r *queryResolver) FrameInfo(ctx context.Context, funcArg string) (*ent.FrameSpec, error) {
	fi, err := r.dbClient.FrameSpec.Query().Where(framespec.Frame(funcArg)).Only(ctx)
	nfe := &ent.NotFoundError{}
	if errors.As(err, &nfe) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return fi, nil
}

// TypeInfo is the resolver for the typeInfo field.
func (r *queryResolver) TypeInfo(ctx context.Context, name string) (*TypeInfo, error) {
	log.Printf("!!! TypeInfo query: %s", name)
	var svcName string
	for serviceName, _ := range r.conf.Targets {
		// TODO(andrei): deal with multiple services
		svcName = serviceName
	}

	var agentURL string
	for _, url := range r.conf.Targets[svcName] {
		agentURL = url
		break
	}

	fields, err := r.getTypeInfoFromDelveAgent(agentURL, name)
	if err != nil {
		return nil, err
	}
	return &TypeInfo{
		Name:            name,
		Fields:          fields,
		FieldsNotLoaded: false,
	}, nil
}

// GetTree is the resolver for the getTree field.
func (r *queryResolver) GetTree(ctx context.Context, colID int, snapID int) (string, error) {
	snap, err := r.stacksFetcher.Fetch(ctx, colID, snapID)
	if err != nil {
		return "", err
	}
	tree := stacks.BuildTree(snap.Snapshot, snap.FramesOfInterest)
	return tree.ToJSON(), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
