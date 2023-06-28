// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package server

type CollectedVar struct {
	Value string  `json:"Value"`
	Links []*Link `json:"Links"`
}

type FieldInfo struct {
	Name     string `json:"Name"`
	Type     string `json:"Type"`
	Embedded bool   `json:"Embedded"`
}

type FrameInfo struct {
	Func string `json:"Func"`
	File string `json:"File"`
	Line int    `json:"Line"`
}

type GoroutineInfo struct {
	ID     int             `json:"ID"`
	Frames []*FrameInfo    `json:"Frames"`
	Vars   []*CollectedVar `json:"Vars"`
}

type GoroutinesGroup struct {
	IDs    []int           `json:"IDs"`
	Frames []*FrameInfo    `json:"Frames"`
	Vars   []*CollectedVar `json:"Vars"`
}

type Link struct {
	SnapshotID  int `json:"SnapshotID"`
	GoroutineID int `json:"GoroutineID"`
	FrameIdx    int `json:"FrameIdx"`
}

type SnapshotInfo struct {
	Raw        []*GoroutineInfo   `json:"Raw"`
	Aggregated []*GoroutinesGroup `json:"Aggregated"`
}

type TypeInfo struct {
	Name            string       `json:"Name"`
	Fields          []*FieldInfo `json:"Fields,omitempty"`
	FieldsNotLoaded bool         `json:"FieldsNotLoaded"`
}

type VarInfo struct {
	Name             string `json:"Name"`
	Type             string `json:"Type"`
	FormalParameter  bool   `json:"FormalParameter"`
	LoclistAvailable bool   `json:"LoclistAvailable"`
}

type VarsAndTypes struct {
	Vars  []*VarInfo  `json:"Vars"`
	Types []*TypeInfo `json:"Types"`
}
