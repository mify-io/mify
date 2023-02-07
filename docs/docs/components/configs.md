---
sidebar_position: 3
---

# Configs

All Mify services have simple generated interface to read static (environment)
and dynamic (Consul) configs.

## Static configs

A static configuration provider is responsible for interacting with any kind of configuration which can't
be changed during application execution. If such configuration changes, it requires application restarting.
Now, this provider supports only ENV configuration.

Let's add such a configuration to `counting-backend` service from our guides.
First, we have to create a new struct that will describe possible values in our new configuration.
For that, navigate to the file with the previously implemented API handler (`go-services/internal/counting-backend/handlers/counter/next/service.go`).

Inside that file, create a new structure:
```go
type CountingAppConf struct {
  IncStep int `envconfig:"INC_STEP" default:"1"`
}
```

This structure contains one int field "IncStep". This field will contain the value loaded from ENV variable
with the name "INC_STEP" (which can be edited in the field tag). Also, we provided a default value "1" for this field. This value
is used if env variable with the name "INC_STEP" doesn't exist.

Now we should access the defined configuration in our code. For doing that, find the line inside API handler where
the counter is increasing (we have added this code in
[Implementing Counter Handler](/docs/guides/create-service/implementing-counter-handler#getting-counter-in-handler) section), and modify this
handler as shown below:
```go
// CounterNextGet - get next number
func (s *CounterNextApiService) CounterNextGet(ctx *core.MifyRequestContext) (openapi.ServiceResponse, error) {
  svcCtx := apputil.GetServiceExtra(ctx.ServiceContext()) // get custom dependencies from context
  currentNumber := svcCtx.Counter

  // NEW CODE
  cfg, err := core.GetStaticConfig[CountingAppConf](ctx)
    if err != nil {
        return openapi.ServiceResponse{}, err
    }

  svcCtx.Counter += cfg.IncStep
  // END OF NEW CODE

  svcCtx.Counter++ // THIS LINE SHOULD BE REMOVED

  return openapi.Response(200, openapi.CounterNextResponse{
    Number: int32(currentNumber),
  }), nil
}
```

Also, another method to access config exists: `MustGetStaticConfig`. This method doesn't return any error but raises
`panic` if any error occurred.

### Building and testing

Now let's test our new code. Since we are using `envconfig`, we should specify the required env variable before
the application is started. For doing that, just run your `counting-backend` like this:
```
$ INC_STEP=3 go run ./cmd/counting-backend
```

Now, you can use curl to validate the result. Don't forget that port in your case could be [different](/docs/guides/create-service/building-and-testing):
```
$ curl 'http://localhost:33767/counter/next'
{"number":0}
$ curl 'http://localhost:33767/counter/next'
{"number":3}
```

Dynamic configs have the same interface as static configs.
