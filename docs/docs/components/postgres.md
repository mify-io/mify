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

### Using without helpers

Here's an example how you can use it, just run queries from pool, which is created on startup:
```go

func (s *PathToApiService) PathToApiGet(ctx *core.MifyRequestContext) (openapi.ServiceResponse, error) {
    rows, err := ctx.Postgres().Query(ctx, "SELECT * FROM table");
    ...
}
```

### Using with sqlc

After adding postgres you will notice `sql-queries/<service_db_name>` directory
inside `go-services`. Create any file with `.sql` extension and refer to
`queries.sql.example` or sqlc [documentation](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html)
for adding queries. After you add them, run `mify generate` to translate them
to go. Generated helpers with lie in `postgres` package, here's an example of
how to use them, assuming you followed sqlc tutorial:

```go
// Call this from service_extra.go
func NewAuthorsStorage(ctx *core.MifyServiceContext) *AuthorsStorage {
    return &AuthorsStorage{
        pool: ctx.Postgres(),
    }
}

func (s *AuthorsStorage) CreateAuthor(
    ctx *core.MifyRequestContext, name string, bio string) (domain.Author, error) {
    dbConn := postgres.New(s.pool)
    tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
    if err != nil {
        return domain.Author{}, err
    }
    defer tx.Rollback(ctx)
    res, err := dbConn.WithTx(tx).CreateAuthor(ctx, postgres.CreateAuthorParams{
        Name: name,
        Bio: bio,
    })
    if err != nil {
        return domain.Author{}, err
    }
    if err := tx.Commit(ctx); err != nil {
        return domain.Author{}, err
    }
    return domain.Author{
        ID: res.ID,
        Name: res.Name,
        Bio: res.Bio,
    }, nil
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

:::info
Ping us in Slack when you ready to use Postgres, we'll set it up for you.
:::

### Connecting to database in Mify Cloud

In mify CLI there is a command for creating ssh session to namespace helper pod:
```
mify cloud ns-shell -e <env>
```
You'll need to provide public ssh key, which will be used for connecting to the pod.
Make sure that your ssh is configured to use this key automatically (e.g. in ssh-agent).

After connecting to the pod, check README in your home directory for the instructions
on how to connect to the database.
