# Primary Key

By default, each model defined in the schema will generate an auto-incremented integer `id` column in the corresponding table. This behavior can be customized using the `primary_key` block.

```hcl
model "Code" {
  default_primary_key = false

  column "type" {
    type = string
    null = false
  }

  column "key" {
    type = string
    null = false
  }

  primary_key {
    columns = ["type", "key"]
  }
}
```

In the example, the `Code` model, which corrsponds to the `codes` table, will have a primary key of `type` and `key`. It is important to note that customizing primary key will affect generated methods, including `Find` and `Delete`. The `Find` method in generated code for the `Code` example will no longer accepts an integer but two strings:

```go
func (q *CodeQuery) Find(typ string, key string) (*Code, error)
```

UUID primary key is common in many application, to support it in queryx:

```hcl
model "Device" {
  default_primary_key = false

  column "id" {
    type = uuid
    null = false
  }

  primary_key {
    columns = ["id"]
  }
}
```
