// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"stacksviz/ent/processsnapshot"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// ProcessSnapshotCreate is the builder for creating a ProcessSnapshot entity.
type ProcessSnapshotCreate struct {
	config
	mutation *ProcessSnapshotMutation
	hooks    []Hook
}

// SetProcessID sets the "process_id" field.
func (psc *ProcessSnapshotCreate) SetProcessID(s string) *ProcessSnapshotCreate {
	psc.mutation.SetProcessID(s)
	return psc
}

// SetSnapshot sets the "snapshot" field.
func (psc *ProcessSnapshotCreate) SetSnapshot(s string) *ProcessSnapshotCreate {
	psc.mutation.SetSnapshot(s)
	return psc
}

// SetFramesOfInterest sets the "frames_of_interest" field.
func (psc *ProcessSnapshotCreate) SetFramesOfInterest(s string) *ProcessSnapshotCreate {
	psc.mutation.SetFramesOfInterest(s)
	return psc
}

// SetNillableFramesOfInterest sets the "frames_of_interest" field if the given value is not nil.
func (psc *ProcessSnapshotCreate) SetNillableFramesOfInterest(s *string) *ProcessSnapshotCreate {
	if s != nil {
		psc.SetFramesOfInterest(*s)
	}
	return psc
}

// SetFlightRecorderData sets the "flight_recorder_data" field.
func (psc *ProcessSnapshotCreate) SetFlightRecorderData(m map[string][]string) *ProcessSnapshotCreate {
	psc.mutation.SetFlightRecorderData(m)
	return psc
}

// Mutation returns the ProcessSnapshotMutation object of the builder.
func (psc *ProcessSnapshotCreate) Mutation() *ProcessSnapshotMutation {
	return psc.mutation
}

// Save creates the ProcessSnapshot in the database.
func (psc *ProcessSnapshotCreate) Save(ctx context.Context) (*ProcessSnapshot, error) {
	return withHooks(ctx, psc.sqlSave, psc.mutation, psc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (psc *ProcessSnapshotCreate) SaveX(ctx context.Context) *ProcessSnapshot {
	v, err := psc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (psc *ProcessSnapshotCreate) Exec(ctx context.Context) error {
	_, err := psc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (psc *ProcessSnapshotCreate) ExecX(ctx context.Context) {
	if err := psc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (psc *ProcessSnapshotCreate) check() error {
	if _, ok := psc.mutation.ProcessID(); !ok {
		return &ValidationError{Name: "process_id", err: errors.New(`ent: missing required field "ProcessSnapshot.process_id"`)}
	}
	if _, ok := psc.mutation.Snapshot(); !ok {
		return &ValidationError{Name: "snapshot", err: errors.New(`ent: missing required field "ProcessSnapshot.snapshot"`)}
	}
	return nil
}

func (psc *ProcessSnapshotCreate) sqlSave(ctx context.Context) (*ProcessSnapshot, error) {
	if err := psc.check(); err != nil {
		return nil, err
	}
	_node, _spec := psc.createSpec()
	if err := sqlgraph.CreateNode(ctx, psc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	psc.mutation.id = &_node.ID
	psc.mutation.done = true
	return _node, nil
}

func (psc *ProcessSnapshotCreate) createSpec() (*ProcessSnapshot, *sqlgraph.CreateSpec) {
	var (
		_node = &ProcessSnapshot{config: psc.config}
		_spec = sqlgraph.NewCreateSpec(processsnapshot.Table, sqlgraph.NewFieldSpec(processsnapshot.FieldID, field.TypeInt))
	)
	if value, ok := psc.mutation.ProcessID(); ok {
		_spec.SetField(processsnapshot.FieldProcessID, field.TypeString, value)
		_node.ProcessID = value
	}
	if value, ok := psc.mutation.Snapshot(); ok {
		_spec.SetField(processsnapshot.FieldSnapshot, field.TypeString, value)
		_node.Snapshot = value
	}
	if value, ok := psc.mutation.FramesOfInterest(); ok {
		_spec.SetField(processsnapshot.FieldFramesOfInterest, field.TypeString, value)
		_node.FramesOfInterest = value
	}
	if value, ok := psc.mutation.FlightRecorderData(); ok {
		_spec.SetField(processsnapshot.FieldFlightRecorderData, field.TypeJSON, value)
		_node.FlightRecorderData = value
	}
	return _node, _spec
}

// ProcessSnapshotCreateBulk is the builder for creating many ProcessSnapshot entities in bulk.
type ProcessSnapshotCreateBulk struct {
	config
	builders []*ProcessSnapshotCreate
}

// Save creates the ProcessSnapshot entities in the database.
func (pscb *ProcessSnapshotCreateBulk) Save(ctx context.Context) ([]*ProcessSnapshot, error) {
	specs := make([]*sqlgraph.CreateSpec, len(pscb.builders))
	nodes := make([]*ProcessSnapshot, len(pscb.builders))
	mutators := make([]Mutator, len(pscb.builders))
	for i := range pscb.builders {
		func(i int, root context.Context) {
			builder := pscb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ProcessSnapshotMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, pscb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, pscb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, pscb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (pscb *ProcessSnapshotCreateBulk) SaveX(ctx context.Context) []*ProcessSnapshot {
	v, err := pscb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pscb *ProcessSnapshotCreateBulk) Exec(ctx context.Context) error {
	_, err := pscb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pscb *ProcessSnapshotCreateBulk) ExecX(ctx context.Context) {
	if err := pscb.Exec(ctx); err != nil {
		panic(err)
	}
}
