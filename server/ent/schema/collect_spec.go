package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CollectSpec describes what data should be collected together with a snapshot:
// what are the frames of interest, and what expressions should be evaluated on
// those frames.
type CollectSpec struct {
	ent.Schema
}

func (CollectSpec) Fields() []ent.Field {
	return []ent.Field{}
}

func (CollectSpec) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("frames", FrameInfo.Type),
	}
}

func (CollectSpec) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
	}
}

type FrameInfo struct {
	ent.Schema
}

func (FrameInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("frame"),
		field.Strings("exprs"),
	}
}

func (FrameInfo) Edges() []ent.Edge {
	return nil
}

func (FrameInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
	}
}
