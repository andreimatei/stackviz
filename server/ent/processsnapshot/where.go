// Code generated by ent, DO NOT EDIT.

package processsnapshot

import (
	"stacksviz/ent/predicate"

	"entgo.io/ent/dialect/sql"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldLTE(FieldID, id))
}

// ProcessID applies equality check predicate on the "process_id" field. It's identical to ProcessIDEQ.
func ProcessID(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldEQ(FieldProcessID, v))
}

// Snapshot applies equality check predicate on the "snapshot" field. It's identical to SnapshotEQ.
func Snapshot(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldEQ(FieldSnapshot, v))
}

// FramesOfInterest applies equality check predicate on the "frames_of_interest" field. It's identical to FramesOfInterestEQ.
func FramesOfInterest(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldEQ(FieldFramesOfInterest, v))
}

// ProcessIDEQ applies the EQ predicate on the "process_id" field.
func ProcessIDEQ(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldEQ(FieldProcessID, v))
}

// ProcessIDNEQ applies the NEQ predicate on the "process_id" field.
func ProcessIDNEQ(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldNEQ(FieldProcessID, v))
}

// ProcessIDIn applies the In predicate on the "process_id" field.
func ProcessIDIn(vs ...string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldIn(FieldProcessID, vs...))
}

// ProcessIDNotIn applies the NotIn predicate on the "process_id" field.
func ProcessIDNotIn(vs ...string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldNotIn(FieldProcessID, vs...))
}

// ProcessIDGT applies the GT predicate on the "process_id" field.
func ProcessIDGT(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldGT(FieldProcessID, v))
}

// ProcessIDGTE applies the GTE predicate on the "process_id" field.
func ProcessIDGTE(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldGTE(FieldProcessID, v))
}

// ProcessIDLT applies the LT predicate on the "process_id" field.
func ProcessIDLT(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldLT(FieldProcessID, v))
}

// ProcessIDLTE applies the LTE predicate on the "process_id" field.
func ProcessIDLTE(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldLTE(FieldProcessID, v))
}

// ProcessIDContains applies the Contains predicate on the "process_id" field.
func ProcessIDContains(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldContains(FieldProcessID, v))
}

// ProcessIDHasPrefix applies the HasPrefix predicate on the "process_id" field.
func ProcessIDHasPrefix(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldHasPrefix(FieldProcessID, v))
}

// ProcessIDHasSuffix applies the HasSuffix predicate on the "process_id" field.
func ProcessIDHasSuffix(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldHasSuffix(FieldProcessID, v))
}

// ProcessIDEqualFold applies the EqualFold predicate on the "process_id" field.
func ProcessIDEqualFold(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldEqualFold(FieldProcessID, v))
}

// ProcessIDContainsFold applies the ContainsFold predicate on the "process_id" field.
func ProcessIDContainsFold(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldContainsFold(FieldProcessID, v))
}

// SnapshotEQ applies the EQ predicate on the "snapshot" field.
func SnapshotEQ(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldEQ(FieldSnapshot, v))
}

// SnapshotNEQ applies the NEQ predicate on the "snapshot" field.
func SnapshotNEQ(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldNEQ(FieldSnapshot, v))
}

// SnapshotIn applies the In predicate on the "snapshot" field.
func SnapshotIn(vs ...string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldIn(FieldSnapshot, vs...))
}

// SnapshotNotIn applies the NotIn predicate on the "snapshot" field.
func SnapshotNotIn(vs ...string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldNotIn(FieldSnapshot, vs...))
}

// SnapshotGT applies the GT predicate on the "snapshot" field.
func SnapshotGT(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldGT(FieldSnapshot, v))
}

// SnapshotGTE applies the GTE predicate on the "snapshot" field.
func SnapshotGTE(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldGTE(FieldSnapshot, v))
}

// SnapshotLT applies the LT predicate on the "snapshot" field.
func SnapshotLT(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldLT(FieldSnapshot, v))
}

// SnapshotLTE applies the LTE predicate on the "snapshot" field.
func SnapshotLTE(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldLTE(FieldSnapshot, v))
}

// SnapshotContains applies the Contains predicate on the "snapshot" field.
func SnapshotContains(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldContains(FieldSnapshot, v))
}

// SnapshotHasPrefix applies the HasPrefix predicate on the "snapshot" field.
func SnapshotHasPrefix(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldHasPrefix(FieldSnapshot, v))
}

// SnapshotHasSuffix applies the HasSuffix predicate on the "snapshot" field.
func SnapshotHasSuffix(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldHasSuffix(FieldSnapshot, v))
}

// SnapshotEqualFold applies the EqualFold predicate on the "snapshot" field.
func SnapshotEqualFold(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldEqualFold(FieldSnapshot, v))
}

// SnapshotContainsFold applies the ContainsFold predicate on the "snapshot" field.
func SnapshotContainsFold(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldContainsFold(FieldSnapshot, v))
}

// FramesOfInterestEQ applies the EQ predicate on the "frames_of_interest" field.
func FramesOfInterestEQ(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldEQ(FieldFramesOfInterest, v))
}

// FramesOfInterestNEQ applies the NEQ predicate on the "frames_of_interest" field.
func FramesOfInterestNEQ(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldNEQ(FieldFramesOfInterest, v))
}

// FramesOfInterestIn applies the In predicate on the "frames_of_interest" field.
func FramesOfInterestIn(vs ...string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldIn(FieldFramesOfInterest, vs...))
}

// FramesOfInterestNotIn applies the NotIn predicate on the "frames_of_interest" field.
func FramesOfInterestNotIn(vs ...string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldNotIn(FieldFramesOfInterest, vs...))
}

// FramesOfInterestGT applies the GT predicate on the "frames_of_interest" field.
func FramesOfInterestGT(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldGT(FieldFramesOfInterest, v))
}

// FramesOfInterestGTE applies the GTE predicate on the "frames_of_interest" field.
func FramesOfInterestGTE(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldGTE(FieldFramesOfInterest, v))
}

// FramesOfInterestLT applies the LT predicate on the "frames_of_interest" field.
func FramesOfInterestLT(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldLT(FieldFramesOfInterest, v))
}

// FramesOfInterestLTE applies the LTE predicate on the "frames_of_interest" field.
func FramesOfInterestLTE(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldLTE(FieldFramesOfInterest, v))
}

// FramesOfInterestContains applies the Contains predicate on the "frames_of_interest" field.
func FramesOfInterestContains(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldContains(FieldFramesOfInterest, v))
}

// FramesOfInterestHasPrefix applies the HasPrefix predicate on the "frames_of_interest" field.
func FramesOfInterestHasPrefix(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldHasPrefix(FieldFramesOfInterest, v))
}

// FramesOfInterestHasSuffix applies the HasSuffix predicate on the "frames_of_interest" field.
func FramesOfInterestHasSuffix(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldHasSuffix(FieldFramesOfInterest, v))
}

// FramesOfInterestIsNil applies the IsNil predicate on the "frames_of_interest" field.
func FramesOfInterestIsNil() predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldIsNull(FieldFramesOfInterest))
}

// FramesOfInterestNotNil applies the NotNil predicate on the "frames_of_interest" field.
func FramesOfInterestNotNil() predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldNotNull(FieldFramesOfInterest))
}

// FramesOfInterestEqualFold applies the EqualFold predicate on the "frames_of_interest" field.
func FramesOfInterestEqualFold(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldEqualFold(FieldFramesOfInterest, v))
}

// FramesOfInterestContainsFold applies the ContainsFold predicate on the "frames_of_interest" field.
func FramesOfInterestContainsFold(v string) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldContainsFold(FieldFramesOfInterest, v))
}

// FlightRecorderDataIsNil applies the IsNil predicate on the "flight_recorder_data" field.
func FlightRecorderDataIsNil() predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldIsNull(FieldFlightRecorderData))
}

// FlightRecorderDataNotNil applies the NotNil predicate on the "flight_recorder_data" field.
func FlightRecorderDataNotNil() predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(sql.FieldNotNull(FieldFlightRecorderData))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ProcessSnapshot) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ProcessSnapshot) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(func(s *sql.Selector) {
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
func Not(p predicate.ProcessSnapshot) predicate.ProcessSnapshot {
	return predicate.ProcessSnapshot(func(s *sql.Selector) {
		p(s.Not())
	})
}
