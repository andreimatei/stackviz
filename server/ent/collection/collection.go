// Code generated by ent, DO NOT EDIT.

package collection

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the collection type in the database.
	Label = "collection"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// EdgeProcessSnapshots holds the string denoting the process_snapshots edge name in mutations.
	EdgeProcessSnapshots = "process_snapshots"
	// Table holds the table name of the collection in the database.
	Table = "collections"
	// ProcessSnapshotsTable is the table that holds the process_snapshots relation/edge.
	ProcessSnapshotsTable = "process_snapshots"
	// ProcessSnapshotsInverseTable is the table name for the ProcessSnapshot entity.
	// It exists in this package in order to avoid circular dependency with the "processsnapshot" package.
	ProcessSnapshotsInverseTable = "process_snapshots"
	// ProcessSnapshotsColumn is the table column denoting the process_snapshots relation/edge.
	ProcessSnapshotsColumn = "collection_process_snapshots"
)

// Columns holds all SQL columns for collection fields.
var Columns = []string{
	FieldID,
	FieldName,
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

var (
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(int) error
)

// OrderOption defines the ordering options for the Collection queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByProcessSnapshotsCount orders the results by process_snapshots count.
func ByProcessSnapshotsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newProcessSnapshotsStep(), opts...)
	}
}

// ByProcessSnapshots orders the results by process_snapshots terms.
func ByProcessSnapshots(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProcessSnapshotsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newProcessSnapshotsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProcessSnapshotsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ProcessSnapshotsTable, ProcessSnapshotsColumn),
	)
}
