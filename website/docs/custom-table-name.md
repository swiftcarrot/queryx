# Custom table name

By default, queryx generates a `table_name` in plural form. For example, a `User` model will have a table named `users`. However, you can customize this behavior using the `table_name` attribute in model block. For example:

```hcl{2}
model "User" {
  table_name = "queryx_users"
}
```

In this example, queryx will generate the table `queryx_users` for the `User` model.
