// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"stacksviz/ent/processsnapshot"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// ProcessSnapshot is the model entity for the ProcessSnapshot schema.
type ProcessSnapshot struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// ProcessID holds the value of the "process_id" field.
	ProcessID string `json:"process_id,omitempty"`
	// Snapshot holds the value of the "snapshot" field.
	Snapshot                     string `json:"snapshot,omitempty"`
	collection_process_snapshots *int
	selectValues                 sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ProcessSnapshot) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case processsnapshot.FieldID:
			values[i] = new(sql.NullInt64)
		case processsnapshot.FieldProcessID, processsnapshot.FieldSnapshot:
			values[i] = new(sql.NullString)
		case processsnapshot.ForeignKeys[0]: // collection_process_snapshots
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ProcessSnapshot fields.
func (ps *ProcessSnapshot) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case processsnapshot.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ps.ID = int(value.Int64)
		case processsnapshot.FieldProcessID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field process_id", values[i])
			} else if value.Valid {
				ps.ProcessID = value.String
			}
		case processsnapshot.FieldSnapshot:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field snapshot", values[i])
			} else if value.Valid {
				ps.Snapshot = value.String
			}
		case processsnapshot.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field collection_process_snapshots", value)
			} else if value.Valid {
				ps.collection_process_snapshots = new(int)
				*ps.collection_process_snapshots = int(value.Int64)
			}
		default:
			ps.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ProcessSnapshot.
// This includes values selected through modifiers, order, etc.
func (ps *ProcessSnapshot) Value(name string) (ent.Value, error) {
	return ps.selectValues.Get(name)
}

// Update returns a builder for updating this ProcessSnapshot.
// Note that you need to call ProcessSnapshot.Unwrap() before calling this method if this ProcessSnapshot
// was returned from a transaction, and the transaction was committed or rolled back.
func (ps *ProcessSnapshot) Update() *ProcessSnapshotUpdateOne {
	return NewProcessSnapshotClient(ps.config).UpdateOne(ps)
}

// Unwrap unwraps the ProcessSnapshot entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ps *ProcessSnapshot) Unwrap() *ProcessSnapshot {
	_tx, ok := ps.config.driver.(*txDriver)
	if !ok {
		panic("ent: ProcessSnapshot is not a transactional entity")
	}
	ps.config.driver = _tx.drv
	return ps
}

// String implements the fmt.Stringer.
func (ps *ProcessSnapshot) String() string {
	var builder strings.Builder
	builder.WriteString("ProcessSnapshot(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ps.ID))
	builder.WriteString("process_id=")
	builder.WriteString(ps.ProcessID)
	builder.WriteString(", ")
	builder.WriteString("snapshot=")
	builder.WriteString(ps.Snapshot)
	builder.WriteByte(')')
	return builder.String()
}

// ProcessSnapshots is a parsable slice of ProcessSnapshot.
type ProcessSnapshots []*ProcessSnapshot