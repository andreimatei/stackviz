// Code generated by ent, DO NOT EDIT.

package collection

import (
	"stacksviz/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Collection {
	return predicate.Collection(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Collection {
	return predicate.Collection(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Collection {
	return predicate.Collection(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Collection {
	return predicate.Collection(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Collection {
	return predicate.Collection(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Collection {
	return predicate.Collection(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Collection {
	return predicate.Collection(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Collection {
	return predicate.Collection(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Collection {
	return predicate.Collection(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Collection {
	return predicate.Collection(sql.FieldEQ(FieldName, v))
}

// CollectSpec applies equality check predicate on the "collect_spec" field. It's identical to CollectSpecEQ.
func CollectSpec(v int) predicate.Collection {
	return predicate.Collection(sql.FieldEQ(FieldCollectSpec, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Collection {
	return predicate.Collection(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Collection {
	return predicate.Collection(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Collection {
	return predicate.Collection(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Collection {
	return predicate.Collection(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Collection {
	return predicate.Collection(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Collection {
	return predicate.Collection(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Collection {
	return predicate.Collection(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Collection {
	return predicate.Collection(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Collection {
	return predicate.Collection(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Collection {
	return predicate.Collection(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Collection {
	return predicate.Collection(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Collection {
	return predicate.Collection(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Collection {
	return predicate.Collection(sql.FieldContainsFold(FieldName, v))
}

// CollectSpecEQ applies the EQ predicate on the "collect_spec" field.
func CollectSpecEQ(v int) predicate.Collection {
	return predicate.Collection(sql.FieldEQ(FieldCollectSpec, v))
}

// CollectSpecNEQ applies the NEQ predicate on the "collect_spec" field.
func CollectSpecNEQ(v int) predicate.Collection {
	return predicate.Collection(sql.FieldNEQ(FieldCollectSpec, v))
}

// CollectSpecIn applies the In predicate on the "collect_spec" field.
func CollectSpecIn(vs ...int) predicate.Collection {
	return predicate.Collection(sql.FieldIn(FieldCollectSpec, vs...))
}

// CollectSpecNotIn applies the NotIn predicate on the "collect_spec" field.
func CollectSpecNotIn(vs ...int) predicate.Collection {
	return predicate.Collection(sql.FieldNotIn(FieldCollectSpec, vs...))
}

// CollectSpecGT applies the GT predicate on the "collect_spec" field.
func CollectSpecGT(v int) predicate.Collection {
	return predicate.Collection(sql.FieldGT(FieldCollectSpec, v))
}

// CollectSpecGTE applies the GTE predicate on the "collect_spec" field.
func CollectSpecGTE(v int) predicate.Collection {
	return predicate.Collection(sql.FieldGTE(FieldCollectSpec, v))
}

// CollectSpecLT applies the LT predicate on the "collect_spec" field.
func CollectSpecLT(v int) predicate.Collection {
	return predicate.Collection(sql.FieldLT(FieldCollectSpec, v))
}

// CollectSpecLTE applies the LTE predicate on the "collect_spec" field.
func CollectSpecLTE(v int) predicate.Collection {
	return predicate.Collection(sql.FieldLTE(FieldCollectSpec, v))
}

// HasProcessSnapshots applies the HasEdge predicate on the "process_snapshots" edge.
func HasProcessSnapshots() predicate.Collection {
	return predicate.Collection(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ProcessSnapshotsTable, ProcessSnapshotsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProcessSnapshotsWith applies the HasEdge predicate on the "process_snapshots" edge with a given conditions (other predicates).
func HasProcessSnapshotsWith(preds ...predicate.ProcessSnapshot) predicate.Collection {
	return predicate.Collection(func(s *sql.Selector) {
		step := newProcessSnapshotsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Collection) predicate.Collection {
	return predicate.Collection(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Collection) predicate.Collection {
	return predicate.Collection(func(s *sql.Selector) {
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
func Not(p predicate.Collection) predicate.Collection {
	return predicate.Collection(func(s *sql.Selector) {
		p(s.Not())
	})
}
