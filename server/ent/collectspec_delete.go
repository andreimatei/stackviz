// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"stacksviz/ent/collectspec"
	"stacksviz/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CollectSpecDelete is the builder for deleting a CollectSpec entity.
type CollectSpecDelete struct {
	config
	hooks    []Hook
	mutation *CollectSpecMutation
}

// Where appends a list predicates to the CollectSpecDelete builder.
func (csd *CollectSpecDelete) Where(ps ...predicate.CollectSpec) *CollectSpecDelete {
	csd.mutation.Where(ps...)
	return csd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (csd *CollectSpecDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, csd.sqlExec, csd.mutation, csd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (csd *CollectSpecDelete) ExecX(ctx context.Context) int {
	n, err := csd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (csd *CollectSpecDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(collectspec.Table, sqlgraph.NewFieldSpec(collectspec.FieldID, field.TypeInt))
	if ps := csd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, csd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	csd.mutation.done = true
	return affected, err
}

// CollectSpecDeleteOne is the builder for deleting a single CollectSpec entity.
type CollectSpecDeleteOne struct {
	csd *CollectSpecDelete
}

// Where appends a list predicates to the CollectSpecDelete builder.
func (csdo *CollectSpecDeleteOne) Where(ps ...predicate.CollectSpec) *CollectSpecDeleteOne {
	csdo.csd.mutation.Where(ps...)
	return csdo
}

// Exec executes the deletion query.
func (csdo *CollectSpecDeleteOne) Exec(ctx context.Context) error {
	n, err := csdo.csd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{collectspec.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (csdo *CollectSpecDeleteOne) ExecX(ctx context.Context) {
	if err := csdo.Exec(ctx); err != nil {
		panic(err)
	}
}