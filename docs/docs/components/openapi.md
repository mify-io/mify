---
sidebar_position: 1
---

# OpenAPI

Mify uses OpenAPI [specification](https://swagger.io/specification/) to
simplify creating API between services and frontends. Check it out to find out
how to use it or use [this
example](https://github.com/OAI/OpenAPI-Specification/blob/main/examples/v3.0/petstore-expanded.yaml),
and here we'll describe additional features that our generation provides on
top.

## Middlewares and ServiceExtra

In `internal/<service-name>/app/service_extra.go` you can see this layout:

```go
// vim: set ft=go:

package app

import (
	"net/http"

	"example.com/namespace/wtest1/go-services/internal/svc1/generated/core"
)

type routerConfig struct {
	Middlewares []func(http.Handler) http.Handler
}

func NewRouterConfig(ctx *core.MifyServiceContext) *routerConfig {
	return &routerConfig {
		Middlewares: []func(http.Handler) http.Handler {
		// Add your middlewares here
		},
	}
}

type ServiceExtra struct {
	// Append your dependencies here
}

func NewServiceExtra(ctx *core.MifyServiceContext) (*ServiceExtra, error) {
	// Here you can do your custom service initialization, prepare dependencies
	extra := &ServiceExtra{
		// Here you can initialize your dependencies
	}
	return extra, nil
}
```

In `NewRouterConfig` you can add your custom middlewares e.g. for auth.

Below it there is `ServiceExtra`. This is your go to place to add dependencies
for components and handlers. The `ServiceExtra` struct will be available in handlers
via `MifyRequestContext`.

## Swagger UI

Swagger UI is available for every backend service on maintenance port (you can
get it from the service startup logs). To be able to call API from it, add your
local and Mify Cloud service url to `schemas/<service-name>/api/api.yaml` to
`servers` field.
