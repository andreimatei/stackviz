```shell
go run ./cmd --resource_root ../client/dist/client --stacks_dir ./datasource
```

We use `ent` as an ORM for SQLite, and we use the `gqlgen` integration for
serving GraphQL queries. Regenerate with:
```shell
go generate ./...
```

The `ent` schema is defined in `ent/schema/*.go`. There's `gqlgen.yml`, which
configures `gqlgen`. This file lists `.graphql` files to use for the schema.
These files are created by hand, except for `ent.graphql` which is created by
`entgql`: the connection between `ent` and `gqlgen` is done in `ent/entc.go`. I
think that tells `ent` to generate `ent.graphql`, which then `gqlgen` uses.

GraphQL queries and mutations defined in the `.graphql` files get a generated
"resolver", which we need to implement, in `collection.resolvers.go`.
