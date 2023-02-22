# vim: set ft=yaml:
version: "2"
sql:
- schema: "{{ .MigrationsDir }}"
  queries: "{{ .QueriesDir }}"
  engine: "postgresql"
  gen:
    go:
      package: "postgres"
      out: "{{ .OutDir }}"
      sql_package: "pgx/v4"
