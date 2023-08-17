# Schema API

## Convention

> **Warning**
> WIP, please refer to test example [here](/internal/integration/postgresql.hcl)

## Time Zone

By default, each database uses the "Local" time zone.

```hcl
database "db" {
  time_zone = "Local" # this is optional
}
```

To specify time zone:

```hcl
database "db" {
  time_zone = "Africa/Lagos"
}
```

## Environment Variable

Queryx provides a convenient feature for reading from environment variables using the built-in `env()` HCL function. It is a common practice for applications to read configuration settings from environment variables in production environments. In the following example, by setting `QUERYX_ENV` to `production`, queryx will automatically read the database connection URL from the `DATABASE_URL` environment variable.

```hcl
database "db" {
  config "development" {
    url = "postgres://postgres:postgres@localhost:5432/blog_development?sslmode=disable"
  }

  config "production" {
    url = env("DATABASE_URL")
  }
}
```

## Database Index

Database index can be declared in schema via the `index` block:

```hcl
model "UserPost" {
  belongs_to "user" {}
  belongs_to "post" {}

  index {
    columns = ["user_id", "post_id"]
    unique  = true
  }
}
```

## Custom table name

By default, queryx generates a `table_name` in plural form. For example, a `User` model will have a table named `users`. However, you can customize this behavior using the `table_name` attribute in model block. For example:

```hcl
model "User" {
  table_name = "queryx_users"
}
```

In this example, queryx will generate the table `queryx_users` for the `User` model.

## Custom primary key

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
