<!-- prettier-ignore-start -->

# Getting Started

## Installation

To easily install the latest version of queryx, open your terminal and run the following command:

```sh
curl -sf https://raw.githubusercontent.com/swiftcarrot/queryx/main/install.sh | sh
```

You can also build queryx from the source following the instructions [here](/docs/build-from-source).

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

```hcl{2,5}
database "db" {
  adapter = "mysql"

  config "development" {
    url = "mysql://root:@127.0.0.1:3306/blog_development"
  }
}
```

Example for SQLite database:

```hcl{2,5}
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

In the sample `schema.hcl`, we can define multiple generators as:

```hcl{2,4}
database "db" {
  generator "client-golang" {}

  generator "client-typescript" {}
}
```

::: tip
You can remove one of the generator declarations to enable only Go or TypeScript code generation.
:::


To generate ORM client methods, simply execute the following command:

```sh
queryx generate
```

Queryx will generate both Golang and TypeScript codes in `db` directory based on the database name. We will cover the basics of CRUD operations (create, read, update, delete) using the queryx generated code. For a more detailed breakdown of the generated methods, please refer to the ORM API section.


### Dependency Installation

Queryx generates codes only depending on third-party database driver implementation.

::: tabs key:lang
== Go
```sh
go get ./...
```
== TypeScript
Queryx relies on `date-fns` for handling dates in TypeScript. Additionally, depending on your database, you will need to install specific packages. Here are the installation commands for different databases:

For PostgreSQL:

```sh
npm install pg @types/pg
```

For MySQL:

```sh
npm install mysql2 @types/node
```

For SQLite:

```sh
npm install better-sqlite3
```
:::

To begin, we instantiate a client object, which serves as the entry point for all interactions with the database.

::: tabs key:lang
== Go
```go
c, err := db.NewClient()
```
== TypeScript
```typescript
let c = newClient();
```
:::

Queryx supports changing database data (insert, update and delete) through a change object. For each model defined in the schema, queryx generates a corresponding change object with setting methods for each field in the model. This ensures the correctness of query and makes it easy to modify data in the database safely.

Create a new post:

::: tabs key:lang
== Go
```go
newPost := c.ChangePost().SetTitle("post title").SetContent("post content")
post, err := c.QueryPost().Create(newPost)
```
== TypeScript
```typescript
let post = await c.queryPost().create({ title: "post title", content: "post content" });
```
:::

Queryx also supports the Active Record pattern, which allows `Update()` or `Delete()` call on the returned queryx record.

Update the post:

::: tabs key:lang
== Go
```go
err := post.Update(c.ChangePost().SetTitle("new post title"))
```
== TypeScript
```typescript
await post.update({ title: "new post title" });
```
:::

Delete the post:

::: tabs key:lang
== Go
```go
err := post.Delete()
```
== TypeScript
```typescript
await post.delete();
```
:::

Queryx also supports update and delete by query.

Update all posts by title:

::: tabs key:lang
== Go
```go
updated, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).UpdateAll(c.ChangePost().SetTitle("new post title"))
```
== TypeScript
```typescript
let updated = await c.queryPost().where(c.postTitle.eq("post title")).updateAll({ title: "new post title" });
```
:::


Delete all posts by title:

::: tabs key:lang
== Go
```go
deleted, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).DeleteAll()
```
== TypeScript
```typescript
let deleted = c.queryPost().where(c.postTitle.eq("post title")).deleteAll();
```
:::


To retrieve data from the database in Queryx using the primary key:

::: tabs key:lang
== Go
```go
post, err := c.QueryPost().Find(1)
```
== TypeScript
```typescript
let post = await c.queryPost().find(1);
```
:::

Retrieve all posts by title:

::: tabs key:lang
== Go
```go
posts, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).All()
```
== TypeScript
```typescript
let posts = await c.queryPost().where(c.postTitle.eq("post title")).all();
```
:::

Retrieve the first post from query:

::: tabs key:lang
== Go
```go
post, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).First()
```
== TypeScript
```typescript
let post = await c.queryPost().where(c.postTitle.eq("post title")).first();
```
:::


<!-- prettier-ignore-end -->
