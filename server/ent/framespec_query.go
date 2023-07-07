// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"
	"stacksviz/ent/collectspec"
	"stacksviz/ent/framespec"
	"stacksviz/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// FrameSpecQuery is the builder for querying FrameSpec entities.
type FrameSpecQuery struct {
	config
	ctx                *QueryContext
	order              []framespec.OrderOption
	inters             []Interceptor
	predicates         []predicate.FrameSpec
	withCollectSpecRef *CollectSpecQuery
	modifiers          []func(*sql.Selector)
	loadTotal          []func(context.Context, []*FrameSpec) error
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the FrameSpecQuery builder.
func (fsq *FrameSpecQuery) Where(ps ...predicate.FrameSpec) *FrameSpecQuery {
	fsq.predicates = append(fsq.predicates, ps...)
	return fsq
}

// Limit the number of records to be returned by this query.
func (fsq *FrameSpecQuery) Limit(limit int) *FrameSpecQuery {
	fsq.ctx.Limit = &limit
	return fsq
}

// Offset to start from.
func (fsq *FrameSpecQuery) Offset(offset int) *FrameSpecQuery {
	fsq.ctx.Offset = &offset
	return fsq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (fsq *FrameSpecQuery) Unique(unique bool) *FrameSpecQuery {
	fsq.ctx.Unique = &unique
	return fsq
}

// Order specifies how the records should be ordered.
func (fsq *FrameSpecQuery) Order(o ...framespec.OrderOption) *FrameSpecQuery {
	fsq.order = append(fsq.order, o...)
	return fsq
}

// QueryCollectSpecRef chains the current query on the "collect_spec_ref" edge.
func (fsq *FrameSpecQuery) QueryCollectSpecRef() *CollectSpecQuery {
	query := (&CollectSpecClient{config: fsq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fsq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := fsq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(framespec.Table, framespec.FieldID, selector),
			sqlgraph.To(collectspec.Table, collectspec.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, framespec.CollectSpecRefTable, framespec.CollectSpecRefColumn),
		)
		fromU = sqlgraph.SetNeighbors(fsq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first FrameSpec entity from the query.
// Returns a *NotFoundError when no FrameSpec was found.
func (fsq *FrameSpecQuery) First(ctx context.Context) (*FrameSpec, error) {
	nodes, err := fsq.Limit(1).All(setContextOp(ctx, fsq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{framespec.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (fsq *FrameSpecQuery) FirstX(ctx context.Context) *FrameSpec {
	node, err := fsq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first FrameSpec ID from the query.
// Returns a *NotFoundError when no FrameSpec ID was found.
func (fsq *FrameSpecQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = fsq.Limit(1).IDs(setContextOp(ctx, fsq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{framespec.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (fsq *FrameSpecQuery) FirstIDX(ctx context.Context) int {
	id, err := fsq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single FrameSpec entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one FrameSpec entity is found.
// Returns a *NotFoundError when no FrameSpec entities are found.
func (fsq *FrameSpecQuery) Only(ctx context.Context) (*FrameSpec, error) {
	nodes, err := fsq.Limit(2).All(setContextOp(ctx, fsq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{framespec.Label}
	default:
		return nil, &NotSingularError{framespec.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (fsq *FrameSpecQuery) OnlyX(ctx context.Context) *FrameSpec {
	node, err := fsq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only FrameSpec ID in the query.
// Returns a *NotSingularError when more than one FrameSpec ID is found.
// Returns a *NotFoundError when no entities are found.
func (fsq *FrameSpecQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = fsq.Limit(2).IDs(setContextOp(ctx, fsq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{framespec.Label}
	default:
		err = &NotSingularError{framespec.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (fsq *FrameSpecQuery) OnlyIDX(ctx context.Context) int {
	id, err := fsq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of FrameSpecs.
func (fsq *FrameSpecQuery) All(ctx context.Context) ([]*FrameSpec, error) {
	ctx = setContextOp(ctx, fsq.ctx, "All")
	if err := fsq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*FrameSpec, *FrameSpecQuery]()
	return withInterceptors[[]*FrameSpec](ctx, fsq, qr, fsq.inters)
}

// AllX is like All, but panics if an error occurs.
func (fsq *FrameSpecQuery) AllX(ctx context.Context) []*FrameSpec {
	nodes, err := fsq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of FrameSpec IDs.
func (fsq *FrameSpecQuery) IDs(ctx context.Context) (ids []int, err error) {
	if fsq.ctx.Unique == nil && fsq.path != nil {
		fsq.Unique(true)
	}
	ctx = setContextOp(ctx, fsq.ctx, "IDs")
	if err = fsq.Select(framespec.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (fsq *FrameSpecQuery) IDsX(ctx context.Context) []int {
	ids, err := fsq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (fsq *FrameSpecQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, fsq.ctx, "Count")
	if err := fsq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, fsq, querierCount[*FrameSpecQuery](), fsq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (fsq *FrameSpecQuery) CountX(ctx context.Context) int {
	count, err := fsq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (fsq *FrameSpecQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, fsq.ctx, "Exist")
	switch _, err := fsq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (fsq *FrameSpecQuery) ExistX(ctx context.Context) bool {
	exist, err := fsq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the FrameSpecQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (fsq *FrameSpecQuery) Clone() *FrameSpecQuery {
	if fsq == nil {
		return nil
	}
	return &FrameSpecQuery{
		config:             fsq.config,
		ctx:                fsq.ctx.Clone(),
		order:              append([]framespec.OrderOption{}, fsq.order...),
		inters:             append([]Interceptor{}, fsq.inters...),
		predicates:         append([]predicate.FrameSpec{}, fsq.predicates...),
		withCollectSpecRef: fsq.withCollectSpecRef.Clone(),
		// clone intermediate query.
		sql:  fsq.sql.Clone(),
		path: fsq.path,
	}
}

// WithCollectSpecRef tells the query-builder to eager-load the nodes that are connected to
// the "collect_spec_ref" edge. The optional arguments are used to configure the query builder of the edge.
func (fsq *FrameSpecQuery) WithCollectSpecRef(opts ...func(*CollectSpecQuery)) *FrameSpecQuery {
	query := (&CollectSpecClient{config: fsq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	fsq.withCollectSpecRef = query
	return fsq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Frame string `json:"frame,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.FrameSpec.Query().
//		GroupBy(framespec.FieldFrame).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (fsq *FrameSpecQuery) GroupBy(field string, fields ...string) *FrameSpecGroupBy {
	fsq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &FrameSpecGroupBy{build: fsq}
	grbuild.flds = &fsq.ctx.Fields
	grbuild.label = framespec.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Frame string `json:"frame,omitempty"`
//	}
//
//	client.FrameSpec.Query().
//		Select(framespec.FieldFrame).
//		Scan(ctx, &v)
func (fsq *FrameSpecQuery) Select(fields ...string) *FrameSpecSelect {
	fsq.ctx.Fields = append(fsq.ctx.Fields, fields...)
	sbuild := &FrameSpecSelect{FrameSpecQuery: fsq}
	sbuild.label = framespec.Label
	sbuild.flds, sbuild.scan = &fsq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a FrameSpecSelect configured with the given aggregations.
func (fsq *FrameSpecQuery) Aggregate(fns ...AggregateFunc) *FrameSpecSelect {
	return fsq.Select().Aggregate(fns...)
}

func (fsq *FrameSpecQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range fsq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, fsq); err != nil {
				return err
			}
		}
	}
	for _, f := range fsq.ctx.Fields {
		if !framespec.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if fsq.path != nil {
		prev, err := fsq.path(ctx)
		if err != nil {
			return err
		}
		fsq.sql = prev
	}
	return nil
}

func (fsq *FrameSpecQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*FrameSpec, error) {
	var (
		nodes       = []*FrameSpec{}
		_spec       = fsq.querySpec()
		loadedTypes = [1]bool{
			fsq.withCollectSpecRef != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*FrameSpec).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &FrameSpec{config: fsq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(fsq.modifiers) > 0 {
		_spec.Modifiers = fsq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, fsq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := fsq.withCollectSpecRef; query != nil {
		if err := fsq.loadCollectSpecRef(ctx, query, nodes, nil,
			func(n *FrameSpec, e *CollectSpec) { n.Edges.CollectSpecRef = e }); err != nil {
			return nil, err
		}
	}
	for i := range fsq.loadTotal {
		if err := fsq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (fsq *FrameSpecQuery) loadCollectSpecRef(ctx context.Context, query *CollectSpecQuery, nodes []*FrameSpec, init func(*FrameSpec), assign func(*FrameSpec, *CollectSpec)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*FrameSpec)
	for i := range nodes {
		fk := nodes[i].CollectSpec
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(collectspec.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "collect_spec" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (fsq *FrameSpecQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := fsq.querySpec()
	if len(fsq.modifiers) > 0 {
		_spec.Modifiers = fsq.modifiers
	}
	_spec.Node.Columns = fsq.ctx.Fields
	if len(fsq.ctx.Fields) > 0 {
		_spec.Unique = fsq.ctx.Unique != nil && *fsq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, fsq.driver, _spec)
}

func (fsq *FrameSpecQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(framespec.Table, framespec.Columns, sqlgraph.NewFieldSpec(framespec.FieldID, field.TypeInt))
	_spec.From = fsq.sql
	if unique := fsq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if fsq.path != nil {
		_spec.Unique = true
	}
	if fields := fsq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, framespec.FieldID)
		for i := range fields {
			if fields[i] != framespec.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if fsq.withCollectSpecRef != nil {
			_spec.Node.AddColumnOnce(framespec.FieldCollectSpec)
		}
	}
	if ps := fsq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := fsq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := fsq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := fsq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (fsq *FrameSpecQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(fsq.driver.Dialect())
	t1 := builder.Table(framespec.Table)
	columns := fsq.ctx.Fields
	if len(columns) == 0 {
		columns = framespec.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if fsq.sql != nil {
		selector = fsq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if fsq.ctx.Unique != nil && *fsq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range fsq.predicates {
		p(selector)
	}
	for _, p := range fsq.order {
		p(selector)
	}
	if offset := fsq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := fsq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// FrameSpecGroupBy is the group-by builder for FrameSpec entities.
type FrameSpecGroupBy struct {
	selector
	build *FrameSpecQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (fsgb *FrameSpecGroupBy) Aggregate(fns ...AggregateFunc) *FrameSpecGroupBy {
	fsgb.fns = append(fsgb.fns, fns...)
	return fsgb
}

// Scan applies the selector query and scans the result into the given value.
func (fsgb *FrameSpecGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, fsgb.build.ctx, "GroupBy")
	if err := fsgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*FrameSpecQuery, *FrameSpecGroupBy](ctx, fsgb.build, fsgb, fsgb.build.inters, v)
}

func (fsgb *FrameSpecGroupBy) sqlScan(ctx context.Context, root *FrameSpecQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(fsgb.fns))
	for _, fn := range fsgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*fsgb.flds)+len(fsgb.fns))
		for _, f := range *fsgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*fsgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := fsgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// FrameSpecSelect is the builder for selecting fields of FrameSpec entities.
type FrameSpecSelect struct {
	*FrameSpecQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (fss *FrameSpecSelect) Aggregate(fns ...AggregateFunc) *FrameSpecSelect {
	fss.fns = append(fss.fns, fns...)
	return fss
}

// Scan applies the selector query and scans the result into the given value.
func (fss *FrameSpecSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, fss.ctx, "Select")
	if err := fss.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*FrameSpecQuery, *FrameSpecSelect](ctx, fss.FrameSpecQuery, fss, fss.inters, v)
}

func (fss *FrameSpecSelect) sqlScan(ctx context.Context, root *FrameSpecQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(fss.fns))
	for _, fn := range fss.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*fss.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := fss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
