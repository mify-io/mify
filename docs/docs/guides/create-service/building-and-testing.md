---
sidebar_position: 2
---

# Building and Testing

Now that we have the code, we can test the handler.

Get into `go-services` directory and install dependencies:

```
$ cd go-services
$ go mod tidy
```

Build and run the service:

```
$ go run ./cmd/counting-backend
```

If everything is ok you will see service logs:
```
{"level":"info","@timestamp":"2022-11-13T01:09:13.220317049Z","caller":"app/mify_app.go:72","msg":"Starting...","service_name":"counting-backend","hostname":"your-hostname"}
{"level":"info","@timestamp":"2022-11-13T01:09:13.220364807Z","caller":"app/server.go:37","msg":"starting api server","service_name":"counting-backend","hostname":"your-hostname","endpoint":":33767"}
{"level":"info","@timestamp":"2022-11-13T01:09:13.220425122Z","caller":"app/server.go:37","msg":"starting maintenance server","service_name":"counting-backend","hostname":"your-hostname","endpoint":":39275"}
```

You can get the port number from `starting api server` log message, we have `33767` in this example.
Use curl or Postman to test your handler:

```
$ curl 'http://localhost:33767/counter/next'
{"number":0}
$ curl 'http://localhost:33767/counter/next'
{"number":1}
```

It works!

