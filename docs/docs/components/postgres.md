---
sidebar_position: 4
---

# Postgres

## How to use

If your service need to use Postgres as storage you need to update
`schemas/<service-name>/service.mify.yaml` and add:

```
postgres:
  enabled: true
  ```

After that you will be able to get `pgxpool.Pool` from `MifyServiceContext`
`Postgres()` method. You can use it to directly make queries or use in some library.

Here's an example how you can use it:
```go

func (s *PathToApiService) PathToApiGet(ctx *core.MifyRequestContext) (openapi.ServiceResponse, error) {
    rows, err := ctx.Postgres().Query(ctx, "SELECT * FROM table");
    ...
}
```

### Migrations

Mify includes a way to apply migrations for a database via dbmate.

To create new migration, run:
```
mify tool migrate <service-name> new <migration-name>
```

Then you'll see migration file at `go-services/migrations/<service-name>/<date>-<migration-name>.sql`
with this content:

```sql
-- migrate:up


-- migrate:down
```

In first block you can add forward migration query, like CREATE TABLE, in down
block you should add corresponding rollback query.

### Testing locally

First you need to start postgres in docker, here's the example docker-compose file:

```yaml
version: '3'
services:
  postgres:
    image: postgres:13
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: passwd
      POSTGRES_DB: <service-name>
    volumes:
      - ~/.cache/mify-db:/var/lib/postgresql/data
    ports:
      - 5432:5432
```

`mify tool migrate` and local config assumes these credentials to connect to database.
After starting postgres you can run `mify tool migrate <service-name> up` to apply migrations.
To rollback run `mify tool migrate <service-name> down`.

### Deploying to Mify Cloud

After you created all migrations and tested them locally, when you run `mify
cloud deploy <service-name>` it will automatically apply all migrations before
deploying new version of the service.

_NOTE:_ Ping us in Slack when you ready to use Postgres, we'll set it up for you.

### Connecting to database in Mify Cloud

In mify CLI there is a command for creating ssh session to namespace helper pod:
```
mify cloud ns-shell -e <env>
```
You'll need to provide public ssh key, which will be used for connecting to the pod.
Make sure that your ssh is configured to use this key automatically (e.g. in ssh-agent).

After connecting to the pod, check README in your home directory for the instructions
on how to connect to the database.
