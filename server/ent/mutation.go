// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"stacksviz/ent/collection"
	"stacksviz/ent/predicate"
	"stacksviz/ent/processsnapshot"
	"sync"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeCollection      = "Collection"
	TypeProcessSnapshot = "ProcessSnapshot"
)

// CollectionMutation represents an operation that mutates the Collection nodes in the graph.
type CollectionMutation struct {
	config
	op                       Op
	typ                      string
	id                       *int
	name                     *string
	clearedFields            map[string]struct{}
	process_snapshots        map[int]struct{}
	removedprocess_snapshots map[int]struct{}
	clearedprocess_snapshots bool
	done                     bool
	oldValue                 func(context.Context) (*Collection, error)
	predicates               []predicate.Collection
}

var _ ent.Mutation = (*CollectionMutation)(nil)

// collectionOption allows management of the mutation configuration using functional options.
type collectionOption func(*CollectionMutation)

// newCollectionMutation creates new mutation for the Collection entity.
func newCollectionMutation(c config, op Op, opts ...collectionOption) *CollectionMutation {
	m := &CollectionMutation{
		config:        c,
		op:            op,
		typ:           TypeCollection,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withCollectionID sets the ID field of the mutation.
func withCollectionID(id int) collectionOption {
	return func(m *CollectionMutation) {
		var (
			err   error
			once  sync.Once
			value *Collection
		)
		m.oldValue = func(ctx context.Context) (*Collection, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Collection.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withCollection sets the old Collection of the mutation.
func withCollection(node *Collection) collectionOption {
	return func(m *CollectionMutation) {
		m.oldValue = func(context.Context) (*Collection, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m CollectionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m CollectionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *CollectionMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *CollectionMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Collection.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *CollectionMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *CollectionMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Collection entity.
// If the Collection object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CollectionMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *CollectionMutation) ResetName() {
	m.name = nil
}

// AddProcessSnapshotIDs adds the "process_snapshots" edge to the ProcessSnapshot entity by ids.
func (m *CollectionMutation) AddProcessSnapshotIDs(ids ...int) {
	if m.process_snapshots == nil {
		m.process_snapshots = make(map[int]struct{})
	}
	for i := range ids {
		m.process_snapshots[ids[i]] = struct{}{}
	}
}

// ClearProcessSnapshots clears the "process_snapshots" edge to the ProcessSnapshot entity.
func (m *CollectionMutation) ClearProcessSnapshots() {
	m.clearedprocess_snapshots = true
}

// ProcessSnapshotsCleared reports if the "process_snapshots" edge to the ProcessSnapshot entity was cleared.
func (m *CollectionMutation) ProcessSnapshotsCleared() bool {
	return m.clearedprocess_snapshots
}

// RemoveProcessSnapshotIDs removes the "process_snapshots" edge to the ProcessSnapshot entity by IDs.
func (m *CollectionMutation) RemoveProcessSnapshotIDs(ids ...int) {
	if m.removedprocess_snapshots == nil {
		m.removedprocess_snapshots = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.process_snapshots, ids[i])
		m.removedprocess_snapshots[ids[i]] = struct{}{}
	}
}

// RemovedProcessSnapshots returns the removed IDs of the "process_snapshots" edge to the ProcessSnapshot entity.
func (m *CollectionMutation) RemovedProcessSnapshotsIDs() (ids []int) {
	for id := range m.removedprocess_snapshots {
		ids = append(ids, id)
	}
	return
}

// ProcessSnapshotsIDs returns the "process_snapshots" edge IDs in the mutation.
func (m *CollectionMutation) ProcessSnapshotsIDs() (ids []int) {
	for id := range m.process_snapshots {
		ids = append(ids, id)
	}
	return
}

// ResetProcessSnapshots resets all changes to the "process_snapshots" edge.
func (m *CollectionMutation) ResetProcessSnapshots() {
	m.process_snapshots = nil
	m.clearedprocess_snapshots = false
	m.removedprocess_snapshots = nil
}

// Where appends a list predicates to the CollectionMutation builder.
func (m *CollectionMutation) Where(ps ...predicate.Collection) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the CollectionMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *CollectionMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Collection, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *CollectionMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *CollectionMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Collection).
func (m *CollectionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *CollectionMutation) Fields() []string {
	fields := make([]string, 0, 1)
	if m.name != nil {
		fields = append(fields, collection.FieldName)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *CollectionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case collection.FieldName:
		return m.Name()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *CollectionMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case collection.FieldName:
		return m.OldName(ctx)
	}
	return nil, fmt.Errorf("unknown Collection field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *CollectionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case collection.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	}
	return fmt.Errorf("unknown Collection field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *CollectionMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *CollectionMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *CollectionMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Collection numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *CollectionMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *CollectionMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *CollectionMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Collection nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *CollectionMutation) ResetField(name string) error {
	switch name {
	case collection.FieldName:
		m.ResetName()
		return nil
	}
	return fmt.Errorf("unknown Collection field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *CollectionMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.process_snapshots != nil {
		edges = append(edges, collection.EdgeProcessSnapshots)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *CollectionMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case collection.EdgeProcessSnapshots:
		ids := make([]ent.Value, 0, len(m.process_snapshots))
		for id := range m.process_snapshots {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *CollectionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedprocess_snapshots != nil {
		edges = append(edges, collection.EdgeProcessSnapshots)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *CollectionMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case collection.EdgeProcessSnapshots:
		ids := make([]ent.Value, 0, len(m.removedprocess_snapshots))
		for id := range m.removedprocess_snapshots {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *CollectionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedprocess_snapshots {
		edges = append(edges, collection.EdgeProcessSnapshots)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *CollectionMutation) EdgeCleared(name string) bool {
	switch name {
	case collection.EdgeProcessSnapshots:
		return m.clearedprocess_snapshots
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *CollectionMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Collection unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *CollectionMutation) ResetEdge(name string) error {
	switch name {
	case collection.EdgeProcessSnapshots:
		m.ResetProcessSnapshots()
		return nil
	}
	return fmt.Errorf("unknown Collection edge %s", name)
}

// ProcessSnapshotMutation represents an operation that mutates the ProcessSnapshot nodes in the graph.
type ProcessSnapshotMutation struct {
	config
	op                 Op
	typ                string
	id                 *int
	process_id         *string
	snapshot           *string
	frames_of_interest *string
	clearedFields      map[string]struct{}
	done               bool
	oldValue           func(context.Context) (*ProcessSnapshot, error)
	predicates         []predicate.ProcessSnapshot
}

var _ ent.Mutation = (*ProcessSnapshotMutation)(nil)

// processsnapshotOption allows management of the mutation configuration using functional options.
type processsnapshotOption func(*ProcessSnapshotMutation)

// newProcessSnapshotMutation creates new mutation for the ProcessSnapshot entity.
func newProcessSnapshotMutation(c config, op Op, opts ...processsnapshotOption) *ProcessSnapshotMutation {
	m := &ProcessSnapshotMutation{
		config:        c,
		op:            op,
		typ:           TypeProcessSnapshot,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withProcessSnapshotID sets the ID field of the mutation.
func withProcessSnapshotID(id int) processsnapshotOption {
	return func(m *ProcessSnapshotMutation) {
		var (
			err   error
			once  sync.Once
			value *ProcessSnapshot
		)
		m.oldValue = func(ctx context.Context) (*ProcessSnapshot, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().ProcessSnapshot.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withProcessSnapshot sets the old ProcessSnapshot of the mutation.
func withProcessSnapshot(node *ProcessSnapshot) processsnapshotOption {
	return func(m *ProcessSnapshotMutation) {
		m.oldValue = func(context.Context) (*ProcessSnapshot, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ProcessSnapshotMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ProcessSnapshotMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ProcessSnapshotMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ProcessSnapshotMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().ProcessSnapshot.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetProcessID sets the "process_id" field.
func (m *ProcessSnapshotMutation) SetProcessID(s string) {
	m.process_id = &s
}

// ProcessID returns the value of the "process_id" field in the mutation.
func (m *ProcessSnapshotMutation) ProcessID() (r string, exists bool) {
	v := m.process_id
	if v == nil {
		return
	}
	return *v, true
}

// OldProcessID returns the old "process_id" field's value of the ProcessSnapshot entity.
// If the ProcessSnapshot object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProcessSnapshotMutation) OldProcessID(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldProcessID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldProcessID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldProcessID: %w", err)
	}
	return oldValue.ProcessID, nil
}

// ResetProcessID resets all changes to the "process_id" field.
func (m *ProcessSnapshotMutation) ResetProcessID() {
	m.process_id = nil
}

// SetSnapshot sets the "snapshot" field.
func (m *ProcessSnapshotMutation) SetSnapshot(s string) {
	m.snapshot = &s
}

// Snapshot returns the value of the "snapshot" field in the mutation.
func (m *ProcessSnapshotMutation) Snapshot() (r string, exists bool) {
	v := m.snapshot
	if v == nil {
		return
	}
	return *v, true
}

// OldSnapshot returns the old "snapshot" field's value of the ProcessSnapshot entity.
// If the ProcessSnapshot object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProcessSnapshotMutation) OldSnapshot(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSnapshot is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSnapshot requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSnapshot: %w", err)
	}
	return oldValue.Snapshot, nil
}

// ResetSnapshot resets all changes to the "snapshot" field.
func (m *ProcessSnapshotMutation) ResetSnapshot() {
	m.snapshot = nil
}

// SetFramesOfInterest sets the "frames_of_interest" field.
func (m *ProcessSnapshotMutation) SetFramesOfInterest(s string) {
	m.frames_of_interest = &s
}

// FramesOfInterest returns the value of the "frames_of_interest" field in the mutation.
func (m *ProcessSnapshotMutation) FramesOfInterest() (r string, exists bool) {
	v := m.frames_of_interest
	if v == nil {
		return
	}
	return *v, true
}

// OldFramesOfInterest returns the old "frames_of_interest" field's value of the ProcessSnapshot entity.
// If the ProcessSnapshot object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProcessSnapshotMutation) OldFramesOfInterest(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldFramesOfInterest is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldFramesOfInterest requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldFramesOfInterest: %w", err)
	}
	return oldValue.FramesOfInterest, nil
}

// ClearFramesOfInterest clears the value of the "frames_of_interest" field.
func (m *ProcessSnapshotMutation) ClearFramesOfInterest() {
	m.frames_of_interest = nil
	m.clearedFields[processsnapshot.FieldFramesOfInterest] = struct{}{}
}

// FramesOfInterestCleared returns if the "frames_of_interest" field was cleared in this mutation.
func (m *ProcessSnapshotMutation) FramesOfInterestCleared() bool {
	_, ok := m.clearedFields[processsnapshot.FieldFramesOfInterest]
	return ok
}

// ResetFramesOfInterest resets all changes to the "frames_of_interest" field.
func (m *ProcessSnapshotMutation) ResetFramesOfInterest() {
	m.frames_of_interest = nil
	delete(m.clearedFields, processsnapshot.FieldFramesOfInterest)
}

// Where appends a list predicates to the ProcessSnapshotMutation builder.
func (m *ProcessSnapshotMutation) Where(ps ...predicate.ProcessSnapshot) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ProcessSnapshotMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ProcessSnapshotMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.ProcessSnapshot, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ProcessSnapshotMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ProcessSnapshotMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (ProcessSnapshot).
func (m *ProcessSnapshotMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ProcessSnapshotMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.process_id != nil {
		fields = append(fields, processsnapshot.FieldProcessID)
	}
	if m.snapshot != nil {
		fields = append(fields, processsnapshot.FieldSnapshot)
	}
	if m.frames_of_interest != nil {
		fields = append(fields, processsnapshot.FieldFramesOfInterest)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ProcessSnapshotMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case processsnapshot.FieldProcessID:
		return m.ProcessID()
	case processsnapshot.FieldSnapshot:
		return m.Snapshot()
	case processsnapshot.FieldFramesOfInterest:
		return m.FramesOfInterest()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ProcessSnapshotMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case processsnapshot.FieldProcessID:
		return m.OldProcessID(ctx)
	case processsnapshot.FieldSnapshot:
		return m.OldSnapshot(ctx)
	case processsnapshot.FieldFramesOfInterest:
		return m.OldFramesOfInterest(ctx)
	}
	return nil, fmt.Errorf("unknown ProcessSnapshot field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ProcessSnapshotMutation) SetField(name string, value ent.Value) error {
	switch name {
	case processsnapshot.FieldProcessID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetProcessID(v)
		return nil
	case processsnapshot.FieldSnapshot:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSnapshot(v)
		return nil
	case processsnapshot.FieldFramesOfInterest:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFramesOfInterest(v)
		return nil
	}
	return fmt.Errorf("unknown ProcessSnapshot field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ProcessSnapshotMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ProcessSnapshotMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ProcessSnapshotMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown ProcessSnapshot numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ProcessSnapshotMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(processsnapshot.FieldFramesOfInterest) {
		fields = append(fields, processsnapshot.FieldFramesOfInterest)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ProcessSnapshotMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ProcessSnapshotMutation) ClearField(name string) error {
	switch name {
	case processsnapshot.FieldFramesOfInterest:
		m.ClearFramesOfInterest()
		return nil
	}
	return fmt.Errorf("unknown ProcessSnapshot nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ProcessSnapshotMutation) ResetField(name string) error {
	switch name {
	case processsnapshot.FieldProcessID:
		m.ResetProcessID()
		return nil
	case processsnapshot.FieldSnapshot:
		m.ResetSnapshot()
		return nil
	case processsnapshot.FieldFramesOfInterest:
		m.ResetFramesOfInterest()
		return nil
	}
	return fmt.Errorf("unknown ProcessSnapshot field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ProcessSnapshotMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ProcessSnapshotMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ProcessSnapshotMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ProcessSnapshotMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ProcessSnapshotMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ProcessSnapshotMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ProcessSnapshotMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown ProcessSnapshot unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ProcessSnapshotMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown ProcessSnapshot edge %s", name)
}
