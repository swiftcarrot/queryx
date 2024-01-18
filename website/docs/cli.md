# CLI

By default, the queryx cli will read from `schema.hcl` in the current directory. To use an alternative schema file, you can specify the file path using the `--schema` flag:

```sh
queryx format --schema db.hcl
```

All available commands:

- `queryx db:create`: create the database
- `queryx db:drop`: drop the database
- `queryx db:migrate`: generate migration files and run pending migrations
- `queryx db:migrate:generate`: generate migration files
  <!-- - `queryx db:migrate:status`: list status of each migration -->
  <!-- - `queryx db:rollback`: rollback database -->
  <!-- - `queryx db:version`: print database migration version -->
- `queryx format`: format schema file with HCL formatter
- `queryx generate`: generate code based on schema
- `queryx version`: print current installed queryx version
