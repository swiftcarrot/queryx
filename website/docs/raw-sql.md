# Raw SQL

Queryx provides direct support for executing raw SQL queries through the following methods.

## Query

`Query` returns all rows from the database.

::: tabs key:lang
== Go

```go
type User struct {
  Name string `db:"user_name"`
}
var users []User
err := c.Query("select name as user_name from users where id in (?)", []int64{1, 2}).Scan(&users)
```

== TypeScript

```typescript
let users = await c.query<{ user_name: string }>(
  "select name as user_name from users where id in (?)",
  [1, 2],
);
```

:::

## Query One

`QueryOne` returns at most a single row from the database.

::: tabs key:lang
== Go

```go
var user struct {
  ID int64 `db:"user_id"`
}
err = c.QueryOne("select id as user_id from users where id = ?", 1).Scan(&user)
```

== TypeScript

```typescript
let user = await c.queryOne<{ user_id: number }>(
  "select id as user_id from users where id = ?",
  1,
);
```

:::

## Exec

`Exec` for SQL statement that don't return data.

::: tabs key:lang
== Go

```go
rowsAffected, err := c.Exec("update users set name = ? where id = ?", "test1", 1)
```

== TypeScript

```typescript
let rowAffected = await c.exec(
  "update users set name = ? where id = ?",
  "test1",
  1,
);
```

:::
