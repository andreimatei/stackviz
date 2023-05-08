Run the frontend with
```shell
ng serve
```

The data requests (i.e. the `/GetData` route) will be proxied to
`localhost:7410` courtesy of the `proxy.conf.json` configuration. The backend is
expected to be running.