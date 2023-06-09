// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func (cs *CollectSpec) Frames(ctx context.Context) (result []*FrameSpec, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = cs.NamedFrames(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = cs.Edges.FramesOrErr()
	}
	if IsNotLoaded(err) {
		result, err = cs.QueryFrames().All(ctx)
	}
	return result, err
}

func (c *Collection) ProcessSnapshots(ctx context.Context) (result []*ProcessSnapshot, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = c.NamedProcessSnapshots(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = c.Edges.ProcessSnapshotsOrErr()
	}
	if IsNotLoaded(err) {
		result, err = c.QueryProcessSnapshots().All(ctx)
	}
	return result, err
}

func (fs *FrameSpec) ParentCollection(ctx context.Context) (*CollectSpec, error) {
	result, err := fs.Edges.ParentCollectionOrErr()
	if IsNotLoaded(err) {
		result, err = fs.QueryParentCollection().Only(ctx)
	}
	return result, err
}
