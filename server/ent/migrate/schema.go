// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// CollectionsColumns holds the columns for the "collections" table.
	CollectionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
	}
	// CollectionsTable holds the schema information for the "collections" table.
	CollectionsTable = &schema.Table{
		Name:       "collections",
		Columns:    CollectionsColumns,
		PrimaryKey: []*schema.Column{CollectionsColumns[0]},
	}
	// ProcessSnapshotsColumns holds the columns for the "process_snapshots" table.
	ProcessSnapshotsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "process_id", Type: field.TypeString},
		{Name: "snapshot", Type: field.TypeString, Size: 2147483647},
		{Name: "collection_process_snapshots", Type: field.TypeInt, Nullable: true},
	}
	// ProcessSnapshotsTable holds the schema information for the "process_snapshots" table.
	ProcessSnapshotsTable = &schema.Table{
		Name:       "process_snapshots",
		Columns:    ProcessSnapshotsColumns,
		PrimaryKey: []*schema.Column{ProcessSnapshotsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "process_snapshots_collections_process_snapshots",
				Columns:    []*schema.Column{ProcessSnapshotsColumns[3]},
				RefColumns: []*schema.Column{CollectionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CollectionsTable,
		ProcessSnapshotsTable,
	}
)

func init() {
	ProcessSnapshotsTable.ForeignKeys[0].RefTable = CollectionsTable
}
