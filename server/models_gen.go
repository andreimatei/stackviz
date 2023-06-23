// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package server

type FieldInfo struct {
	Name     string `json:"Name"`
	Type     string `json:"Type"`
	Embedded bool   `json:"Embedded"`
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
