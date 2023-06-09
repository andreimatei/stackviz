// Code generated by ent, DO NOT EDIT.

package framespec

import (
	"stacksviz/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldLTE(FieldID, id))
}

// Frame applies equality check predicate on the "frame" field. It's identical to FrameEQ.
func Frame(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldEQ(FieldFrame, v))
}

// CollectSpecID applies equality check predicate on the "collect_spec_id" field. It's identical to CollectSpecIDEQ.
func CollectSpecID(v int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldEQ(FieldCollectSpecID, v))
}

// FrameEQ applies the EQ predicate on the "frame" field.
func FrameEQ(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldEQ(FieldFrame, v))
}

// FrameNEQ applies the NEQ predicate on the "frame" field.
func FrameNEQ(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldNEQ(FieldFrame, v))
}

// FrameIn applies the In predicate on the "frame" field.
func FrameIn(vs ...string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldIn(FieldFrame, vs...))
}

// FrameNotIn applies the NotIn predicate on the "frame" field.
func FrameNotIn(vs ...string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldNotIn(FieldFrame, vs...))
}

// FrameGT applies the GT predicate on the "frame" field.
func FrameGT(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldGT(FieldFrame, v))
}

// FrameGTE applies the GTE predicate on the "frame" field.
func FrameGTE(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldGTE(FieldFrame, v))
}

// FrameLT applies the LT predicate on the "frame" field.
func FrameLT(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldLT(FieldFrame, v))
}

// FrameLTE applies the LTE predicate on the "frame" field.
func FrameLTE(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldLTE(FieldFrame, v))
}

// FrameContains applies the Contains predicate on the "frame" field.
func FrameContains(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldContains(FieldFrame, v))
}

// FrameHasPrefix applies the HasPrefix predicate on the "frame" field.
func FrameHasPrefix(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldHasPrefix(FieldFrame, v))
}

// FrameHasSuffix applies the HasSuffix predicate on the "frame" field.
func FrameHasSuffix(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldHasSuffix(FieldFrame, v))
}

// FrameEqualFold applies the EqualFold predicate on the "frame" field.
func FrameEqualFold(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldEqualFold(FieldFrame, v))
}

// FrameContainsFold applies the ContainsFold predicate on the "frame" field.
func FrameContainsFold(v string) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldContainsFold(FieldFrame, v))
}

// CollectSpecIDEQ applies the EQ predicate on the "collect_spec_id" field.
func CollectSpecIDEQ(v int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldEQ(FieldCollectSpecID, v))
}

// CollectSpecIDNEQ applies the NEQ predicate on the "collect_spec_id" field.
func CollectSpecIDNEQ(v int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldNEQ(FieldCollectSpecID, v))
}

// CollectSpecIDIn applies the In predicate on the "collect_spec_id" field.
func CollectSpecIDIn(vs ...int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldIn(FieldCollectSpecID, vs...))
}

// CollectSpecIDNotIn applies the NotIn predicate on the "collect_spec_id" field.
func CollectSpecIDNotIn(vs ...int) predicate.FrameSpec {
	return predicate.FrameSpec(sql.FieldNotIn(FieldCollectSpecID, vs...))
}

// HasParentCollection applies the HasEdge predicate on the "parentCollection" edge.
func HasParentCollection() predicate.FrameSpec {
	return predicate.FrameSpec(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ParentCollectionTable, ParentCollectionColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasParentCollectionWith applies the HasEdge predicate on the "parentCollection" edge with a given conditions (other predicates).
func HasParentCollectionWith(preds ...predicate.CollectSpec) predicate.FrameSpec {
	return predicate.FrameSpec(func(s *sql.Selector) {
		step := newParentCollectionStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.FrameSpec) predicate.FrameSpec {
	return predicate.FrameSpec(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.FrameSpec) predicate.FrameSpec {
	return predicate.FrameSpec(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.FrameSpec) predicate.FrameSpec {
	return predicate.FrameSpec(func(s *sql.Selector) {
		p(s.Not())
	})
}
