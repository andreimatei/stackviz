// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"
	"stacksviz/ent/collectspec"
	"stacksviz/ent/framespec"
	"stacksviz/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CollectSpecQuery is the builder for querying CollectSpec entities.
type CollectSpecQuery struct {
	config
	ctx             *QueryContext
	order           []collectspec.OrderOption
	inters          []Interceptor
	predicates      []predicate.CollectSpec
	withFrames      *FrameSpecQuery
	modifiers       []func(*sql.Selector)
	loadTotal       []func(context.Context, []*CollectSpec) error
	withNamedFrames map[string]*FrameSpecQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CollectSpecQuery builder.
func (csq *CollectSpecQuery) Where(ps ...predicate.CollectSpec) *CollectSpecQuery {
	csq.predicates = append(csq.predicates, ps...)
	return csq
}

// Limit the number of records to be returned by this query.
func (csq *CollectSpecQuery) Limit(limit int) *CollectSpecQuery {
	csq.ctx.Limit = &limit
	return csq
}

// Offset to start from.
func (csq *CollectSpecQuery) Offset(offset int) *CollectSpecQuery {
	csq.ctx.Offset = &offset
	return csq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (csq *CollectSpecQuery) Unique(unique bool) *CollectSpecQuery {
	csq.ctx.Unique = &unique
	return csq
}

// Order specifies how the records should be ordered.
func (csq *CollectSpecQuery) Order(o ...collectspec.OrderOption) *CollectSpecQuery {
	csq.order = append(csq.order, o...)
	return csq
}

// QueryFrames chains the current query on the "frames" edge.
func (csq *CollectSpecQuery) QueryFrames() *FrameSpecQuery {
	query := (&FrameSpecClient{config: csq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := csq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := csq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(collectspec.Table, collectspec.FieldID, selector),
			sqlgraph.To(framespec.Table, framespec.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, collectspec.FramesTable, collectspec.FramesColumn),
		)
		fromU = sqlgraph.SetNeighbors(csq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CollectSpec entity from the query.
// Returns a *NotFoundError when no CollectSpec was found.
func (csq *CollectSpecQuery) First(ctx context.Context) (*CollectSpec, error) {
	nodes, err := csq.Limit(1).All(setContextOp(ctx, csq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{collectspec.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (csq *CollectSpecQuery) FirstX(ctx context.Context) *CollectSpec {
	node, err := csq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first CollectSpec ID from the query.
// Returns a *NotFoundError when no CollectSpec ID was found.
func (csq *CollectSpecQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = csq.Limit(1).IDs(setContextOp(ctx, csq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{collectspec.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (csq *CollectSpecQuery) FirstIDX(ctx context.Context) int {
	id, err := csq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single CollectSpec entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one CollectSpec entity is found.
// Returns a *NotFoundError when no CollectSpec entities are found.
func (csq *CollectSpecQuery) Only(ctx context.Context) (*CollectSpec, error) {
	nodes, err := csq.Limit(2).All(setContextOp(ctx, csq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{collectspec.Label}
	default:
		return nil, &NotSingularError{collectspec.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (csq *CollectSpecQuery) OnlyX(ctx context.Context) *CollectSpec {
	node, err := csq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only CollectSpec ID in the query.
// Returns a *NotSingularError when more than one CollectSpec ID is found.
// Returns a *NotFoundError when no entities are found.
func (csq *CollectSpecQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = csq.Limit(2).IDs(setContextOp(ctx, csq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{collectspec.Label}
	default:
		err = &NotSingularError{collectspec.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (csq *CollectSpecQuery) OnlyIDX(ctx context.Context) int {
	id, err := csq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CollectSpecs.
func (csq *CollectSpecQuery) All(ctx context.Context) ([]*CollectSpec, error) {
	ctx = setContextOp(ctx, csq.ctx, "All")
	if err := csq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*CollectSpec, *CollectSpecQuery]()
	return withInterceptors[[]*CollectSpec](ctx, csq, qr, csq.inters)
}

// AllX is like All, but panics if an error occurs.
func (csq *CollectSpecQuery) AllX(ctx context.Context) []*CollectSpec {
	nodes, err := csq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of CollectSpec IDs.
func (csq *CollectSpecQuery) IDs(ctx context.Context) (ids []int, err error) {
	if csq.ctx.Unique == nil && csq.path != nil {
		csq.Unique(true)
	}
	ctx = setContextOp(ctx, csq.ctx, "IDs")
	if err = csq.Select(collectspec.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (csq *CollectSpecQuery) IDsX(ctx context.Context) []int {
	ids, err := csq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (csq *CollectSpecQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, csq.ctx, "Count")
	if err := csq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, csq, querierCount[*CollectSpecQuery](), csq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (csq *CollectSpecQuery) CountX(ctx context.Context) int {
	count, err := csq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (csq *CollectSpecQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, csq.ctx, "Exist")
	switch _, err := csq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (csq *CollectSpecQuery) ExistX(ctx context.Context) bool {
	exist, err := csq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CollectSpecQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (csq *CollectSpecQuery) Clone() *CollectSpecQuery {
	if csq == nil {
		return nil
	}
	return &CollectSpecQuery{
		config:     csq.config,
		ctx:        csq.ctx.Clone(),
		order:      append([]collectspec.OrderOption{}, csq.order...),
		inters:     append([]Interceptor{}, csq.inters...),
		predicates: append([]predicate.CollectSpec{}, csq.predicates...),
		withFrames: csq.withFrames.Clone(),
		// clone intermediate query.
		sql:  csq.sql.Clone(),
		path: csq.path,
	}
}

// WithFrames tells the query-builder to eager-load the nodes that are connected to
// the "frames" edge. The optional arguments are used to configure the query builder of the edge.
func (csq *CollectSpecQuery) WithFrames(opts ...func(*FrameSpecQuery)) *CollectSpecQuery {
	query := (&FrameSpecClient{config: csq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	csq.withFrames = query
	return csq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (csq *CollectSpecQuery) GroupBy(field string, fields ...string) *CollectSpecGroupBy {
	csq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &CollectSpecGroupBy{build: csq}
	grbuild.flds = &csq.ctx.Fields
	grbuild.label = collectspec.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (csq *CollectSpecQuery) Select(fields ...string) *CollectSpecSelect {
	csq.ctx.Fields = append(csq.ctx.Fields, fields...)
	sbuild := &CollectSpecSelect{CollectSpecQuery: csq}
	sbuild.label = collectspec.Label
	sbuild.flds, sbuild.scan = &csq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a CollectSpecSelect configured with the given aggregations.
func (csq *CollectSpecQuery) Aggregate(fns ...AggregateFunc) *CollectSpecSelect {
	return csq.Select().Aggregate(fns...)
}

func (csq *CollectSpecQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range csq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, csq); err != nil {
				return err
			}
		}
	}
	for _, f := range csq.ctx.Fields {
		if !collectspec.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if csq.path != nil {
		prev, err := csq.path(ctx)
		if err != nil {
			return err
		}
		csq.sql = prev
	}
	return nil
}

func (csq *CollectSpecQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*CollectSpec, error) {
	var (
		nodes       = []*CollectSpec{}
		_spec       = csq.querySpec()
		loadedTypes = [1]bool{
			csq.withFrames != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*CollectSpec).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &CollectSpec{config: csq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(csq.modifiers) > 0 {
		_spec.Modifiers = csq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, csq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := csq.withFrames; query != nil {
		if err := csq.loadFrames(ctx, query, nodes,
			func(n *CollectSpec) { n.Edges.Frames = []*FrameSpec{} },
			func(n *CollectSpec, e *FrameSpec) { n.Edges.Frames = append(n.Edges.Frames, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range csq.withNamedFrames {
		if err := csq.loadFrames(ctx, query, nodes,
			func(n *CollectSpec) { n.appendNamedFrames(name) },
			func(n *CollectSpec, e *FrameSpec) { n.appendNamedFrames(name, e) }); err != nil {
			return nil, err
		}
	}
	for i := range csq.loadTotal {
		if err := csq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (csq *CollectSpecQuery) loadFrames(ctx context.Context, query *FrameSpecQuery, nodes []*CollectSpec, init func(*CollectSpec), assign func(*CollectSpec, *FrameSpec)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*CollectSpec)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(framespec.FieldCollectSpecID)
	}
	query.Where(predicate.FrameSpec(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(collectspec.FramesColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.CollectSpecID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "collect_spec_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (csq *CollectSpecQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := csq.querySpec()
	if len(csq.modifiers) > 0 {
		_spec.Modifiers = csq.modifiers
	}
	_spec.Node.Columns = csq.ctx.Fields
	if len(csq.ctx.Fields) > 0 {
		_spec.Unique = csq.ctx.Unique != nil && *csq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, csq.driver, _spec)
}

func (csq *CollectSpecQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(collectspec.Table, collectspec.Columns, sqlgraph.NewFieldSpec(collectspec.FieldID, field.TypeInt))
	_spec.From = csq.sql
	if unique := csq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if csq.path != nil {
		_spec.Unique = true
	}
	if fields := csq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, collectspec.FieldID)
		for i := range fields {
			if fields[i] != collectspec.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := csq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := csq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := csq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := csq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (csq *CollectSpecQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(csq.driver.Dialect())
	t1 := builder.Table(collectspec.Table)
	columns := csq.ctx.Fields
	if len(columns) == 0 {
		columns = collectspec.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if csq.sql != nil {
		selector = csq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if csq.ctx.Unique != nil && *csq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range csq.predicates {
		p(selector)
	}
	for _, p := range csq.order {
		p(selector)
	}
	if offset := csq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := csq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// WithNamedFrames tells the query-builder to eager-load the nodes that are connected to the "frames"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (csq *CollectSpecQuery) WithNamedFrames(name string, opts ...func(*FrameSpecQuery)) *CollectSpecQuery {
	query := (&FrameSpecClient{config: csq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if csq.withNamedFrames == nil {
		csq.withNamedFrames = make(map[string]*FrameSpecQuery)
	}
	csq.withNamedFrames[name] = query
	return csq
}

// CollectSpecGroupBy is the group-by builder for CollectSpec entities.
type CollectSpecGroupBy struct {
	selector
	build *CollectSpecQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (csgb *CollectSpecGroupBy) Aggregate(fns ...AggregateFunc) *CollectSpecGroupBy {
	csgb.fns = append(csgb.fns, fns...)
	return csgb
}

// Scan applies the selector query and scans the result into the given value.
func (csgb *CollectSpecGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, csgb.build.ctx, "GroupBy")
	if err := csgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CollectSpecQuery, *CollectSpecGroupBy](ctx, csgb.build, csgb, csgb.build.inters, v)
}

func (csgb *CollectSpecGroupBy) sqlScan(ctx context.Context, root *CollectSpecQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(csgb.fns))
	for _, fn := range csgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*csgb.flds)+len(csgb.fns))
		for _, f := range *csgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*csgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := csgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// CollectSpecSelect is the builder for selecting fields of CollectSpec entities.
type CollectSpecSelect struct {
	*CollectSpecQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (css *CollectSpecSelect) Aggregate(fns ...AggregateFunc) *CollectSpecSelect {
	css.fns = append(css.fns, fns...)
	return css
}

// Scan applies the selector query and scans the result into the given value.
func (css *CollectSpecSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, css.ctx, "Select")
	if err := css.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CollectSpecQuery, *CollectSpecSelect](ctx, css.CollectSpecQuery, css, css.inters, v)
}

func (css *CollectSpecSelect) sqlScan(ctx context.Context, root *CollectSpecQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(css.fns))
	for _, fn := range css.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*css.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := css.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
