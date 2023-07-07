package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
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
		edge.To("frames", FrameSpec.Type).
			Annotations(entsql.OnDelete(entsql.Cascade), entgql.RelayConnection()),
	}
}

func (CollectSpec) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
	}
}

type FrameSpec struct {
	ent.Schema
}

func (FrameSpec) Fields() []ent.Field {
	return []ent.Field{
		field.String("frame"),
		field.Int("collect_spec").Comment("The parent collection spec"),
		// collect_expressions is the list of expressions to evaluate whenever this
		// frame is encountered when taking a snapshot. Each expression is a string
		// that is passed to Delve's `eval` function.
		field.Strings("collect_expressions"),
		// flight_recorder_events is the list of events to record
		// whenever this function begins execution.
		// Each event is a JSON string with the following format:
		// {
		//   "expr": "<the expression to be eval()'ed>",
		//   "key_expr": "<the expression to be used as the key>",
		// }
		field.Strings("flight_recorder_events"),
	}
}

func (FrameSpec) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("collect_spec_ref", CollectSpec.Type).Ref("frames").
			Field("collect_spec").Required().Unique(),
	}
}

func (FrameSpec) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
		entgql.RelayConnection(),
	}
}
