// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"stacksviz/ent/collection"
	"stacksviz/ent/processsnapshot"

	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
)

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (c *CollectionQuery) CollectFields(ctx context.Context, satisfies ...string) (*CollectionQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return c, nil
	}
	if err := c.collectField(ctx, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *CollectionQuery) collectField(ctx context.Context, opCtx *graphql.OperationContext, collected graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	var (
		unknownSeen    bool
		fieldSeen      = make(map[string]struct{}, len(collection.Columns))
		selectedFields = []string{collection.FieldID}
	)
	for _, field := range graphql.CollectFields(opCtx, collected.Selections, satisfies) {
		switch field.Name {
		case "processSnapshots":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = (&ProcessSnapshotClient{config: c.config}).Query()
			)
			if err := query.collectField(ctx, opCtx, field, path, mayAddCondition(satisfies, processsnapshotImplementors)...); err != nil {
				return err
			}
			c.WithNamedProcessSnapshots(alias, func(wq *ProcessSnapshotQuery) {
				*wq = *query
			})
		case "name":
			if _, ok := fieldSeen[collection.FieldName]; !ok {
				selectedFields = append(selectedFields, collection.FieldName)
				fieldSeen[collection.FieldName] = struct{}{}
			}
		case "id":
		case "__typename":
		default:
			unknownSeen = true
		}
	}
	if !unknownSeen {
		c.Select(selectedFields...)
	}
	return nil
}

type collectionPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []CollectionPaginateOption
}

func newCollectionPaginateArgs(rv map[string]any) *collectionPaginateArgs {
	args := &collectionPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	return args
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ps *ProcessSnapshotQuery) CollectFields(ctx context.Context, satisfies ...string) (*ProcessSnapshotQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return ps, nil
	}
	if err := ps.collectField(ctx, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return ps, nil
}

func (ps *ProcessSnapshotQuery) collectField(ctx context.Context, opCtx *graphql.OperationContext, collected graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	var (
		unknownSeen    bool
		fieldSeen      = make(map[string]struct{}, len(processsnapshot.Columns))
		selectedFields = []string{processsnapshot.FieldID}
	)
	for _, field := range graphql.CollectFields(opCtx, collected.Selections, satisfies) {
		switch field.Name {
		case "processID":
			if _, ok := fieldSeen[processsnapshot.FieldProcessID]; !ok {
				selectedFields = append(selectedFields, processsnapshot.FieldProcessID)
				fieldSeen[processsnapshot.FieldProcessID] = struct{}{}
			}
		case "snapshot":
			if _, ok := fieldSeen[processsnapshot.FieldSnapshot]; !ok {
				selectedFields = append(selectedFields, processsnapshot.FieldSnapshot)
				fieldSeen[processsnapshot.FieldSnapshot] = struct{}{}
			}
		case "framesOfInterest":
			if _, ok := fieldSeen[processsnapshot.FieldFramesOfInterest]; !ok {
				selectedFields = append(selectedFields, processsnapshot.FieldFramesOfInterest)
				fieldSeen[processsnapshot.FieldFramesOfInterest] = struct{}{}
			}
		case "id":
		case "__typename":
		default:
			unknownSeen = true
		}
	}
	if !unknownSeen {
		ps.Select(selectedFields...)
	}
	return nil
}

type processsnapshotPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []ProcessSnapshotPaginateOption
}

func newProcessSnapshotPaginateArgs(rv map[string]any) *processsnapshotPaginateArgs {
	args := &processsnapshotPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	return args
}

const (
	afterField     = "after"
	firstField     = "first"
	beforeField    = "before"
	lastField      = "last"
	orderByField   = "orderBy"
	directionField = "direction"
	fieldField     = "field"
	whereField     = "where"
)

func fieldArgs(ctx context.Context, whereInput any, path ...string) map[string]any {
	field := collectedField(ctx, path...)
	if field == nil || field.Arguments == nil {
		return nil
	}
	oc := graphql.GetOperationContext(ctx)
	args := field.ArgumentMap(oc.Variables)
	return unmarshalArgs(ctx, whereInput, args)
}

// unmarshalArgs allows extracting the field arguments from their raw representation.
func unmarshalArgs(ctx context.Context, whereInput any, args map[string]any) map[string]any {
	for _, k := range []string{firstField, lastField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		i, err := graphql.UnmarshalInt(v)
		if err == nil {
			args[k] = &i
		}
	}
	for _, k := range []string{beforeField, afterField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		c := &Cursor{}
		if c.UnmarshalGQL(v) == nil {
			args[k] = c
		}
	}
	if v, ok := args[whereField]; ok && whereInput != nil {
		if err := graphql.UnmarshalInputFromContext(ctx, v, whereInput); err == nil {
			args[whereField] = whereInput
		}
	}

	return args
}

func limitRows(partitionBy string, limit int, orderBy ...sql.Querier) func(s *sql.Selector) {
	return func(s *sql.Selector) {
		d := sql.Dialect(s.Dialect())
		s.SetDistinct(false)
		with := d.With("src_query").
			As(s.Clone()).
			With("limited_query").
			As(
				d.Select("*").
					AppendSelectExprAs(
						sql.RowNumber().PartitionBy(partitionBy).OrderExpr(orderBy...),
						"row_number",
					).
					From(d.Table("src_query")),
			)
		t := d.Table("limited_query").As(s.TableName())
		*s = *d.Select(s.UnqualifiedColumns()...).
			From(t).
			Where(sql.LTE(t.C("row_number"), limit)).
			Prefix(with)
	}
}

// mayAddCondition appends another type condition to the satisfies list
// if it does not exist in the list.
func mayAddCondition(satisfies []string, typeCond []string) []string {
Cond:
	for _, c := range typeCond {
		for _, s := range satisfies {
			if c == s {
				continue Cond
			}
		}
		satisfies = append(satisfies, c)
	}
	return satisfies
}
