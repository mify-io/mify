---
sidebar_position: 0
---

# Using configuration in your application

All backend services generated with mify contain code required for accessing application configs.

## Working with static configuration

A static configuration provider is responsible for interacting with any kind of configuration which can't
be changed during application execution. If such configuration changes, it requires application restarting.
Now, this provider supports only ENV configuration.

Let's add such a configuration to `counting-backend` service.
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
is used if env variable with the name "INC_STEP" will not exist.

Now we should access the defined configuration in our code. For doing that, find the line inside API handler where
the counter is increasing (we have added this code in
[Implementing Counter Handler](/create-service/implementing-counter-handler) section), and modify this
handler as shown below:
```go
// CounterNextGet - get next number
func (s *CounterNextApiService) CounterNextGet(ctx *core.MifyRequestContext) (openapi.ServiceResponse, error) {
  svcCtx := ctx.ServiceExtra().(*app.ServiceExtra) // get custom dependencies from context
  currentNumber := svcCtx.Counter

  // NEW CODE
  cfg, err := ctx.StaticConfig().Get(CountingAppConf{})
  if err != nil {
    return openapi.ServiceResponse{}, err
  }

  svcCtx.Counter += cfg.(CountingAppConf).IncStep
  // END OF NEW CODE

  svcCtx.Counter++ // THIS LINE SHOULD BE REMOVED

  return openapi.Response(200, openapi.CounterNextResponse{
    Number: int32(currentNumber),
  }), nil
}
```

So, for accessing static configuration we are using ```ctx.StaticConfig()```.

## Building and testing

Now we should check the new code. Since we are using `envconfig`, we should specify the required env variable before
the application is started. For doing that, just run your `counting-backend` like this:
```
export INC_STEP=3; go run ./cmd/counting-backend
```

Now, you can use curl to validate the result:
```
$ curl 'http://localhost:33767/counter/next'
{"number":0}
$ curl 'http://localhost:33767/counter/next'
{"number":3}
```
