package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// ProcessSnapshot holds the schema definition for the ProcessSnapshot entity.
type ProcessSnapshot struct {
	ent.Schema
}

// Fields of the ProcessSnapshot.
func (ProcessSnapshot) Fields() []ent.Field {
	return []ent.Field{
		field.String("process_id"),
		field.Text("snapshot"),
		field.Text("frames_of_interest").Optional(),
	}
}

// Edges of the ProcessSnapshot.
func (ProcessSnapshot) Edges() []ent.Edge {
	return nil
}

func (ProcessSnapshot) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
	}
}
