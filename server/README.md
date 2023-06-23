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

---

# Running everything:

```shell
./cockroach start-single-node --insecure --logtostderr
./workload workload run kv --read-percent=100 --concurrency=10
[sudo sh -c " echo 0 > /proc/sys/kernel/yama/ptrace_scope"]
./dlv attach --accept-multiclient `pidof cockroach` --listen=127.0.0.1:45689 --headless --log

[/home/andrei/src/github.com/andreimatei/delve-agent]
go run ./cmd/agent.go

[/home/andrei/src/github.com/andreimatei/stackviz/server]
go run ./cmd/main.go --resource_root ../client/dist/client --stacks_dir ./datasource

[/home/andrei/src/github.com/andreimatei/stackviz/client]
ng serve
```

I'm using my fork of traceviz, starlark-go and panicparse.
