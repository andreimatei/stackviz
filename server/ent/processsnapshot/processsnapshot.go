// Code generated by ent, DO NOT EDIT.

package processsnapshot

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the processsnapshot type in the database.
	Label = "process_snapshot"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldProcessID holds the string denoting the process_id field in the database.
	FieldProcessID = "process_id"
	// FieldSnapshot holds the string denoting the snapshot field in the database.
	FieldSnapshot = "snapshot"
	// FieldFramesOfInterest holds the string denoting the frames_of_interest field in the database.
	FieldFramesOfInterest = "frames_of_interest"
	// Table holds the table name of the processsnapshot in the database.
	Table = "process_snapshots"
)

// Columns holds all SQL columns for processsnapshot fields.
var Columns = []string{
	FieldID,
	FieldProcessID,
	FieldSnapshot,
	FieldFramesOfInterest,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "process_snapshots"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"collection_process_snapshots",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the ProcessSnapshot queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByProcessID orders the results by the process_id field.
func ByProcessID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProcessID, opts...).ToFunc()
}

// BySnapshot orders the results by the snapshot field.
func BySnapshot(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSnapshot, opts...).ToFunc()
}

// ByFramesOfInterest orders the results by the frames_of_interest field.
func ByFramesOfInterest(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFramesOfInterest, opts...).ToFunc()
}
