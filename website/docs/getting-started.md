# Getting Started

## Installation

To easily install the latest version of queryx, open your terminal and run the following command:

```sh
curl -sf https://raw.githubusercontent.com/swiftcarrot/queryx/main/install.sh | sh
```

You can also build queryx from the source following the instructions [here](/docs/build).

After installation, run the following command to validate queryx:

```sh
queryx version
```

This command will output current installed queryx version if installed successfully.

## Create your first schema

Queryx uses [HCL format](https://github.com/hashicorp/hcl) for schema defintion. Create the following sample `schema.hcl` in the current directory:

```hcl
database "db" {
  adapter = "postgresql"

  config "development" {
    url = "postgresql://postgres:postgres@localhost:5432/blog_development?sslmode=disable"
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
    url = "mysql://root:@127.0.0.1:3306/blog_development"
  }
}
```

Example for SQLite database:

```hcl
database "db" {
  adapter = "sqlite"

  config "development" {
    url = "sqlite:blog_development.sqlite3"
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
