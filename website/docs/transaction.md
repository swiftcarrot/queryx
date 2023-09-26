<!-- prettier-ignore-start -->

# Transaction

Queryx also supported type-safe database transactions, making it easy to execute database transactions safely.

Creating a transaction:

::: tabs key:lang
== Go
```go
c, err := db.NewClient()
tx := c.Tx()
```
== TypeScript
```typescript
let c = newClient();
let tx = await c.tx();
```
:::

The queryx transaction object works similarly to the queryx client methods with the exception that it requires an additional commit call to make changes to the database.

::: tabs key:lang
== Go

```go
post, err := tx.QueryPost().Create(tx.ChangPost().SetTitle("post title"))
err := post.Update(tx.ChangePost().SetTitle("new post title"))
if err := tx.Commit() {
  tx.Rollback()
}
```

== TypeScript

```typescript
try {
  let post = await tx.queryPost().create(tx.changPost().setTitle("post title"));
  await post.update(tx.changePost().setTitle("new post title"));
  await tx.Commit();
} catch (err) {
  await tx.Rollback();
}
```

:::

Queryx also supports wrapping transactions within a function block to avoid manual commit and rollback operations. Queryx automatically performs the rollback operation in case of an error returned from the function block.

::: tabs key:lang
== Go
```go
err := c.Transaction(func (tx *db.Tx) error {
  post, err := tx.QueryPost().Create(tx.ChangPost().SetTitle("post title"))
  if err != nil {
    return err
  }
  err := post.Update(tx.ChangePost().SetTitle("new post title"))
  if err != nil {
    return err
  }
  return nil
})
```
== TypeScript
```typescript
await c.transaction(async function (tx: Tx) {
  let post = await tx.queryPost().create(tx.changePost().setTitle("post title"));
  await post.update(tx.changePost().setTitle("new post title"));
});
```
:::

<!-- prettier-ignore-end -->
