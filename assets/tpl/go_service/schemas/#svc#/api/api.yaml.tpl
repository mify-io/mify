openapi: "3.0.0"
info:
  version: 1.0.0
  title: {{.ServiceName}}
  description: Service description
  contact:
    name: Maintainer name
    email: Maintainer email
servers:
  - url: {{.ServiceName}}.company.com
paths:
  /path/to/api:
    get:
      summary: sample handler
      operationId: theOperationId
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
