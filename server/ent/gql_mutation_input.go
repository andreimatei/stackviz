// Code generated by ent, DO NOT EDIT.

package ent

// CreateCollectionInput represents a mutation input for creating collections.
type CreateCollectionInput struct {
	Name               string
	ProcessSnapshotIDs []int
}

// Mutate applies the CreateCollectionInput on the CollectionMutation builder.
func (i *CreateCollectionInput) Mutate(m *CollectionMutation) {
	m.SetName(i.Name)
	if v := i.ProcessSnapshotIDs; len(v) > 0 {
		m.AddProcessSnapshotIDs(v...)
	}
}

// SetInput applies the change-set in the CreateCollectionInput on the CollectionCreate builder.
func (c *CollectionCreate) SetInput(i CreateCollectionInput) *CollectionCreate {
	i.Mutate(c.Mutation())
	return c
}

// CreateProcessSnapshotInput represents a mutation input for creating processsnapshots.
type CreateProcessSnapshotInput struct {
	ProcessID        string
	Snapshot         string
	FramesOfInterest []string
}

// Mutate applies the CreateProcessSnapshotInput on the ProcessSnapshotMutation builder.
func (i *CreateProcessSnapshotInput) Mutate(m *ProcessSnapshotMutation) {
	m.SetProcessID(i.ProcessID)
	m.SetSnapshot(i.Snapshot)
	if v := i.FramesOfInterest; v != nil {
		m.SetFramesOfInterest(v)
	}
}

// SetInput applies the change-set in the CreateProcessSnapshotInput on the ProcessSnapshotCreate builder.
func (c *ProcessSnapshotCreate) SetInput(i CreateProcessSnapshotInput) *ProcessSnapshotCreate {
	i.Mutate(c.Mutation())
	return c
}
