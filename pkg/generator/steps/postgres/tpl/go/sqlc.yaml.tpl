# vim: set ft=yaml:
version: "2"
sql:
- schema: "{{ .Model.MigrationsDir }}"
  queries: "{{ .Model.QueriesDir }}"
  engine: "postgresql"
  gen:
    go:
      package: "postgres"
      out: "{{ .Model.OutDir }}"
      sql_package: "pgx/v4"
