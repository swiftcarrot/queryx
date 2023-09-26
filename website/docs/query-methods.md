# Query Methods

::: warning
The following document is a work in progress.
:::

For each model defined in the schema, queryx generates a corresponding query object.

::: tabs key:lang
== Go

```go
q := c.QueryPost()
```

== TypeScript

```typescript
let q = c.queryPost();
```

:::

### Finder Methods

A query object supports the following find methods:

- Find
- FindBy
- FindBySQL

### Query Methods

Query contruction:

- Where
- Limit
- Offset
- Order
- Joins

Query execution:

- All
- First
- Count
- Exists
- UpdateAll
- DeleteAll

### Record Methods

- Update
- Delete
