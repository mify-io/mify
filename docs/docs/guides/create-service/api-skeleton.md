---
sidebar_position: 0
---

# Generating Skeleton with API

First we need to add service to workspace:
```
$ mify add service counting-backend
```

For this tutorial we're using Go service template which tries to follow this
layout (https://github.com/golang-standards/project-layout). This is not
an official standard project layout, because Go doesn't have one, but it's pretty common.

Now we need to define the API for this service. Open OpenAPI schema file, which
is located at `schemas/counting-backend/api/api.yaml` and you will see this
schema:

```yaml
openapi: "3.0.0"
info:
  version: 1.0.0
  title: counting-backend
  description: Service description
  contact:
    name: Maintainer name
    email: Maintainer email
servers:
  - url: counting-backend.example.com
paths: {}
# Example of a handler, uncomment and remove the above 'paths: {}' line.
# Check Petstore OpenAPI example for more possible options:
# https://github.com/OAI/OpenAPI-Specification/blob/main/examples/v3.0/petstore-expanded.yaml
#
# paths:
#   /path/to/api:
#     get:
#       summary: sample handler
#       responses:
#         '200':
#           description: OK
#           content:
#             application/json:
#               schema:
#                 $ref: '#/components/schemas/PathToApiResponse'
# components:
#   schemas:
#     PathToApiResponse:
#       type: object
#       properties:
#         value:
#           type: string
#       required:
#         - value
```

Let's add our handler definition to it:
```yaml
openapi: "3.0.0"
info:
  version: 1.0.0
  title: counting-backend
  description: Service description
  contact:
    name: Maintainer name
    email: Maintainer email
servers:
  - url: counting-backend.example.com
paths:
  /counter/next:
    get:
      summary: get next number
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/CounterNextResponse'

components:
  schemas:
    CounterNextResponse:
      type: object
      properties:
        number:
          type: integer
      required:
        - number
```

Now run command to apply changed schema:

```
mify generate counting-backend
```

*NOTE: If you remove or change the name of the handler after you ran `mify
generate` you need to delete directory with generated files and then re-run
`mify generate` command. Directory with generated files is located in
`go-services/internal/counting-backend/generated`.*

Now we have our handler and we are ready to implement the logic.
