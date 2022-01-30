---
sidebar_position: 2
---

# Create Your First Project

To create project you need to call:

```
mify init <project-name>
```

This will create workspace where you can add your services.
Go into workpace with `cd <project-name>` and create your first service:
```
mify add service <service-name>
```
You will see that this command generated a service in `go-services` directory.
Then, you can start writing your code immediately after that, here is the structure of generated code:

`go-services/cmd/<service-name>` - main package of service

`go-services/internal/<service-name>/app/service_extra.go` - service entrypoint, there you can add your custom dependencies.

`go-services/internal/<service-name>/handlers/.../service.go` - write your handler logic here.

`schemas/<service-name>/api/api.yaml` - add your handlers here and call `mify generate <service-name>` to regenerate code.
