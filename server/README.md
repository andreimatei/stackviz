```shell
go run ./cmd --resource_root ../client/dist/client --stacks_dir ./datasource
```

We use `ent` as an ORM for SQLite, and we use the `gqlgen` integration for
serving GraphQL queries. Regenerate with:
```shell
go generate ./...
```

The `ent` schema is defined in `ent/schema/collection.go` and
`ent/schema/processsnapshot.go`. There's `gqlgen.yml`, which configures
`gqlgen`. The connection between `ent` and `gqlgen` is done in `ent/entc.go`. I
think that tells `ent` to generate `ent.graphql`, which then `gqlgen` uses.
