// Code generated by ent, DO NOT EDIT.

package ent

// CreateCollectSpecInput represents a mutation input for creating collectspecs.
type CreateCollectSpecInput struct {
	FrameIDs []int
}

// Mutate applies the CreateCollectSpecInput on the CollectSpecMutation builder.
func (i *CreateCollectSpecInput) Mutate(m *CollectSpecMutation) {
	if v := i.FrameIDs; len(v) > 0 {
		m.AddFrameIDs(v...)
	}
}

// SetInput applies the change-set in the CreateCollectSpecInput on the CollectSpecCreate builder.
func (c *CollectSpecCreate) SetInput(i CreateCollectSpecInput) *CollectSpecCreate {
	i.Mutate(c.Mutation())
	return c
}

// CreateCollectionInput represents a mutation input for creating collections.
type CreateCollectionInput struct {
	Name               string
	CollectSpec        int
	ProcessSnapshotIDs []int
}

// Mutate applies the CreateCollectionInput on the CollectionMutation builder.
func (i *CreateCollectionInput) Mutate(m *CollectionMutation) {
	m.SetName(i.Name)
	m.SetCollectSpec(i.CollectSpec)
	if v := i.ProcessSnapshotIDs; len(v) > 0 {
		m.AddProcessSnapshotIDs(v...)
	}
}

// SetInput applies the change-set in the CreateCollectionInput on the CollectionCreate builder.
func (c *CollectionCreate) SetInput(i CreateCollectionInput) *CollectionCreate {
	i.Mutate(c.Mutation())
	return c
}

// CreateFrameSpecInput represents a mutation input for creating framespecs.
type CreateFrameSpecInput struct {
	Frame                string
	CollectExpressions   []string
	FlightRecorderEvents []string
	ParentCollectionID   int
}

// Mutate applies the CreateFrameSpecInput on the FrameSpecMutation builder.
func (i *CreateFrameSpecInput) Mutate(m *FrameSpecMutation) {
	m.SetFrame(i.Frame)
	if v := i.CollectExpressions; v != nil {
		m.SetCollectExpressions(v)
	}
	if v := i.FlightRecorderEvents; v != nil {
		m.SetFlightRecorderEvents(v)
	}
	m.SetParentCollectionID(i.ParentCollectionID)
}

// SetInput applies the change-set in the CreateFrameSpecInput on the FrameSpecCreate builder.
func (c *FrameSpecCreate) SetInput(i CreateFrameSpecInput) *FrameSpecCreate {
	i.Mutate(c.Mutation())
	return c
}

// CreateProcessSnapshotInput represents a mutation input for creating processsnapshots.
type CreateProcessSnapshotInput struct {
	ProcessID        string
	Snapshot         string
	FramesOfInterest *string
}

// Mutate applies the CreateProcessSnapshotInput on the ProcessSnapshotMutation builder.
func (i *CreateProcessSnapshotInput) Mutate(m *ProcessSnapshotMutation) {
	m.SetProcessID(i.ProcessID)
	m.SetSnapshot(i.Snapshot)
	if v := i.FramesOfInterest; v != nil {
		m.SetFramesOfInterest(*v)
	}
}

// SetInput applies the change-set in the CreateProcessSnapshotInput on the ProcessSnapshotCreate builder.
func (c *ProcessSnapshotCreate) SetInput(i CreateProcessSnapshotInput) *ProcessSnapshotCreate {
	i.Mutate(c.Mutation())
	return c
}
