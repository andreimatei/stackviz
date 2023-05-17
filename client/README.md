# Running

Run the frontend with
```shell
ng serve
```

The data requests (i.e. the `/GetData` and `/graphql` routes) will be proxied to
`localhost:7410` courtesy of the `proxy.conf.json` configuration. The backend is
expected to be running.

# Development

We're using `apollo-angular` as a GraphQL client library. We're also using
`graphql-codegen` to generate code based on a schema (which it reads from the
server, through introspection) and a list of queries (see
`src/app/graphql/collection.graphql`). The code generation is run by
```shell
npm run graphql-codegen
```
