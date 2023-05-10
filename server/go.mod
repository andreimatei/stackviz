module stacksviz

go 1.19

require (
	github.com/google/traceviz/server/go v0.0.0-20230428161057-a446b048e906
	github.com/hashicorp/golang-lru/v2 v2.0.2
	github.com/maruel/panicparse/v2 v2.3.1
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/safehtml v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.3.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/google/traceviz/server/go => ../../../google/traceviz/server/go
