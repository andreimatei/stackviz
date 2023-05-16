package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// ProcessSnapshot holds the schema definition for the ProcessSnapshot entity.
type ProcessSnapshot struct {
	ent.Schema
}

// Fields of the ProcessSnapshot.
func (ProcessSnapshot) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("process_id"),
		field.Text("snapshot"),
	}
}

// Edges of the ProcessSnapshot.
func (ProcessSnapshot) Edges() []ent.Edge {
	return nil
}
