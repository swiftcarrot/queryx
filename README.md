# Queryx

[English](README.md) | [中文](README_zh.md)

> **Warning**
> This project is currently in beta (v0), although it has been battled tested in internal projects. Currently, it only supports golang code generation. We plan to release support for TypeScript code generation. Feel free to [open an issue](https://github.com/swiftcarrot/queryx/issues) or [start a discussion](https://github.com/swiftcarrot/queryx/discussions) if you have any questions.

Queryx is schema-first and type-safe ORM with code generation.

[![go report card](https://goreportcard.com/badge/github.com/swiftcarrot/queryx "go report card")](https://goreportcard.com/report/github.com/swiftcarrot/queryx)
[![test status](https://github.com/swiftcarrot/queryx/workflows/integration/badge.svg "test status")](https://github.com/swiftcarrot/queryx/actions)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/swiftcarrot/queryx)

- **Schema First**: Define application models in a queryx schema file, and it can automatically synchronize with the database structure.
- **Type Safe**: Queryx generates friendly and type-safe ORM methods based on the schema, which come with autocomplete support and are free from type-related errors.

This project is heavily inspired by [Active Record](https://guides.rubyonrails.org/active_record_basics.html) and [ent](https://entgo.io/). Database schema management is built upon [Atlas](https://atlasgo.io/).

# Getting Started

## Installation

To easily install the latest version of queryx, open your terminal and run the following command:

```sh
curl -sf https://raw.githubusercontent.com/swiftcarrot/queryx/main/install.sh  | sh
```

You can also build queryx from the source following the instructions [here](/docs/build).

After installation, run the following command to validate queryx:

```sh
queryx version
```

This command will output current installed queryx version if installed successfully.

## Define your first schema

Queryx uses [HCL format](https://github.com/hashicorp/hcl) for schema defintion. Create the following sample `schema.hcl` in the current directory:

```hcl
database "db" {
  adapter = "postgresql"

  config "development" {
    url = "postgres://postgres:postgres@localhost:5432/blog_development?sslmode=disable"
  }

  generator "client-golang" {}

  model "Post" {
    column "title" {
      type = string
    }
    column "content" {
      type = text
    }
  }
}
```

In this sample schema, we create a queryx database `db`, which consists of a model `Post`. `Post` model contains two fields, `title` as `string` type and `content` as `text` type. `string` and `text` are both predefined queryx types. The `db` database is defined as the `postgresql` adapter and the connection config url to the PostgreSQL database is defined through the `config` block.

Queryx also support MySQL and SQLite databases by changing the `adapter` attribute and `config` in `database` block.

Example for MySQL database:

```hcl
database "db" {
  adapter = "mysql"

  config "development" {
    url = "root@tcp(localhost:3306)/queryx_test?parseTime=true"
  }
}
```

Example for SQLite database:

```hcl
database "db" {
  adapter = "sqlite"

  config "test" {
    url = "file:test.sqlite3"
  }
}
```

Run the following command to automatically format the schema file:

```sh
queryx format
```

## Database managment

Run the following command to create the PostgreSQL database, by default, queryx with read from the `development` config block. It can be changed by setting the `QUERYX_ENV` environment variable.

```sh
queryx db:create
```

which works the same as

```sh
QUERYX_ENV=development queryx db:create
```

Once the database is created, queryx can automatically migrate the database to the schema defined in `schema.hcl`:

```sh
queryx db:migrate
```

The `db:migrate` command will initially compare the current state of database to the schema defined in `schema.hcl`. It will generate migrations files in SQL format in the `db/migrations` directory if there are any differences. It will proceed to execute the generated migration files to update the database in line with the schema defined in `schema.hcl`.

## Code generation

In the sample `schema.hcl`, we already defined the generator as

```hcl
generator "client-golang"
```

To generate Golang client methods, simply execute the following command:

```sh
queryx generate
```

Queryx will generate Golang codes in `db` directory based on the database name. We will cover the basics of CRUD operations (create, read, update, delete) using the queryx generated code. For a more detailed breakdown of the generated methods, please refer to the ORM API section.

To begin, we instantiate a client object, which serves as the entry point for all interactions with the database.

```go
c, err := db.NewClient()
```

Queryx supports changing database data (insert, update and delete) through a change object. For each model defined in the schema, queryx generates a corresponding change object with setting methods for each field in the model. This ensures the correctness of query and makes it easy to modify data in the database safely.

Create a new post:

```go
newPost := c.ChangePost().SetTitle("post title").SetContent("post content")
post, err := c.QueryPost().Create(newPost)
```

Queryx also supports the Active Record pattern, which allows `Update()` or `Delete()` call on the returned queryx record.

Update the post:

```go
err := post.Update(c.ChangePost().SetTitle("new post title"))
```

Delete the post:

```go
err := post.Delete()
```

Queryx also supports update and delete by query.

Update all posts by title:

```go
updated, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).UpdateAll(c.ChangePost().SetTitle("new post title"))
```

Delete all posts by title:

```go
deleted, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).DeleteAll()
```

To retrieve data from the database in Queryx using the primary key:

```go
post, err := c.QueryPost().Find(1)
```

Retrieve all posts by title:

```go
posts, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).All()
```

Retrieve the first post from query:

```go
post, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).First()
```

# Association

Queryx supports association definition in the schema file. It also generates corresponding preload query methods to avoid "N+1" query.

## belongs_to

```hcl
model "Post" {
  belongs_to "Author" {
    model_name = "User"
  }
}
```

```go
c.QueryPost().PreloadAuthor().All()
```

## has_one

```hcl
model "User" {
  has_one "account" {}

  column "name" {
    type = string
  }
}

model "Account" {
  belongs_to "user" {}

  column "name" {
    type = string
  }
}
```

![](./docs/has_one.png)

```go
c.QueryUser().PreloadAccount().All()
c.QueryAccount().PreloadUser().All()
```

## has_many

```hcl
model "User" {
  belongs_to "group" {}

  column "name" {
    type = string
  }
}

model "Group" {
  has_many "users" {}

  column "name" {
    type = string
  }
}
```

![](./docs/has_many.png)

```go
c.QueryUser().PreloadGroup().All()
c.QueryGroup().PreloadUsers().All()
```

## has_many through

```hcl
model "User" {
  has_many "user_posts" {}
  has_many "posts" {
    through = "user_posts"
  }
}

model "Post" {
  has_many "user_posts" {}
  has_many "users" {
    through = "user_posts"
  }
}

model "UserPost" {
  belongs_to "user" {}
  belongs_to "post" {}
}
```

![](./docs/has_many_through.png)

```go
c.QueryUser().PreloadPosts().All()
c.QueryPost().PreloadUsers().All()
```

# ORM API

## Query

For each model defined in the schema, queryx generates a corresponding query object.

```go
q := c.QueryPost()
```

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

## Transaction

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

# Data Types

Predefined data types in queryx:

- `integer`:
- `bigint`:
- `string`:
- `text`:
- `boolean`: A true/false value
- `float`:
- `json/jsonb`:
- `uuid`:
- `datetime`: A time and date (2006-01-02 15:04:05)
- `time`: A time without date (2006-01-02)
- `date`: A date without time (15:04:05)

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

# License

Queryx is licensed under Apache 2.0 as found in the LICENSE file.
