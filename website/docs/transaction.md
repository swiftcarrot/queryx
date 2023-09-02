# Transaction

Queryx also supported type-safe database transactions, making it easy to execute database transactions safely.

Creating a transaction:

```go
c, err := db.NewClient()
tx := c.Tx()
```

The queryx transaction object works similarly to the queryx client methods with the exception that it requires an additional commit call to make changes to the database.

```go
post, err := tx.QueryPost().Create(tx.ChangPost().SetTitle("post title"))
err := post.Update(tx.ChangePost().SetTitle("new post title"))
if err := tx.Commit() {
  tx.Rollback()
}
```

