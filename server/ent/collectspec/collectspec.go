// Code generated by ent, DO NOT EDIT.

package collectspec

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the collectspec type in the database.
	Label = "collect_spec"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// EdgeFrames holds the string denoting the frames edge name in mutations.
	EdgeFrames = "frames"
	// Table holds the table name of the collectspec in the database.
	Table = "collect_specs"
	// FramesTable is the table that holds the frames relation/edge.
	FramesTable = "frame_specs"
	// FramesInverseTable is the table name for the FrameSpec entity.
	// It exists in this package in order to avoid circular dependency with the "framespec" package.
	FramesInverseTable = "frame_specs"
	// FramesColumn is the table column denoting the frames relation/edge.
	FramesColumn = "collect_spec"
)

// Columns holds all SQL columns for collectspec fields.
var Columns = []string{
	FieldID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the CollectSpec queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByFramesCount orders the results by frames count.
func ByFramesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newFramesStep(), opts...)
	}
}

// ByFrames orders the results by frames terms.
func ByFrames(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newFramesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newFramesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(FramesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, FramesTable, FramesColumn),
	)
}
