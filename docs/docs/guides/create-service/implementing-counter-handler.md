---
sidebar_position: 1
---

# Implementing Counter Handler

## Storage for counter

Before implementing handler we want to store our current counter somewhere.
For the purpose of this tutorial we would store it in memory, but where should we define it?

Every handler needs some common dependencies, like logger, config, clients, and
for that Mify provides `MifyRequestContext` struct. It also supports custom
dependencies, and this is where we can put our counter, so let's do it!

Open `go-services/internal/counting-backend/app/service_extra.go`,
add `Counter` field to `ServiceExtra` struct (which is already defined in the end of the file),
and you should end up with something like this:
```go
type ServiceExtra struct {
	// Append your dependencies here
	Counter int
}

func NewServiceExtra(ctx *core.MifyServiceContext) (*ServiceExtra, error) {
	// Here you can do your custom service initialization, prepare dependencies
	extra := &ServiceExtra{
		// Here you can initialize your dependencies
		Counter: 0,
	}
	return extra, nil
}
```

## Getting counter in handler

Now we can finally implement the handler (`go-services/internal/counting-backend/handlers/counter/next/service.go`):

```go
// CounterNextGet - get next number
func (s *CounterNextApiService) CounterNextGet(ctx *core.MifyRequestContext) (openapi.ServiceResponse, error) {
	svcCtx := apputil.GetServiceExtra(ctx.ServiceContext()) // get custom dependencies from context
	currentNumber := svcCtx.Counter

	svcCtx.Counter++

	return openapi.Response(200, openapi.CounterNextResponse{
		Number: int32(currentNumber),
	}), nil
}
```

Add import for `apputil.GetServiceExtra` and remove unused ones:
```
import (
	"example.com/namespace/counting-project/go-services/internal/counting-backend/apputil"
	"example.com/namespace/counting-project/go-services/internal/counting-backend/generated/api"
	"example.com/namespace/counting-project/go-services/internal/counting-backend/generated/core"
)
```
