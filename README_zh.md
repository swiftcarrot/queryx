# Queryx

[English](README.md) | [中文](README_zh.md)

> **Warning**
> 该项目目前处于 Beta 版本（v0），虽然已在内部项目中进行了测试。目前，它仅支持生成 golang 代码。我们计划发布支持生成 TypeScript 代码的功能。 如果您有任何问题，请随时 [提出 issue](https://github.com/swiftcarrot/queryx/issues) 或 [展开讨论](https://github.com/swiftcarrot/queryx/discussions)。

> TypeScript 支持已发布！现在可以在项目中使用 TypeScript 来使用 queryx，[查看文档](/docs/generator-typescript-client.md)进行操作。

<img width="400" alt="image" src="https://github.com/swiftcarrot/queryx/assets/1039026/db0f79b9-4bc8-4beb-ab28-7a8de7d8f8e9">

Queryx 是一个基于模式优先、类型安全的 ORM，具有代码生成功能。

[![go report card](https://goreportcard.com/badge/github.com/swiftcarrot/queryx "go report card")](https://goreportcard.com/report/github.com/swiftcarrot/queryx)
[![test status](https://github.com/swiftcarrot/queryx/workflows/integration/badge.svg "test status")](https://github.com/swiftcarrot/queryx/actions)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/swiftcarrot/queryx)

- **模式优先**: 在 queryx 的模式文件中定义应用程序模型，自动与数据库结构同步。
- **类型安全**: Queryx 根据模式生成友好的、类型安全的 ORM 方法，这些方法具有自动完成支持，避免出现与类型相关的错误。

这个项目的灵感来源于 [Active Record](https://guides.rubyonrails.org/active_record_basics.html) 和 [ent](https://entgo.io/). 数据库模式管理是基于 [Atlas](https://atlasgo.io/)。

# 开始

## 安装

要轻松安装最新版本的 Queryx，请打开终端并运行以下命令：

```sh
curl -sf https://raw.githubusercontent.com/swiftcarrot/queryx/main/install.sh  | sh
```

如果您希望自行构建 Queryx，您可以按照以下步骤操作 [这里](/docs/build).

在安装完成后，您可以运行以下命令来验证 Queryx 是否被正确安装：

```sh
queryx version
```

如果 Queryx 成功安装，这个命令会输出当前程序的版本号.

## 定义你的第一个模式

Queryx 使用 [HCL format](https://github.com/hashicorp/hcl) 来定义模式. 在当前目录创建以下示例文件 `schema.hcl`:

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

在这个示例模式中，我们创建了一个名为 db 的 queryx 数据库，其中包含一个名为 Post 的模型。模型 Post 包含两个字段，title 为 string 类型，content 为 text 类型。string 和 text 均为 queryx 预定义的类型。db 数据库被定义为 postgresql 适配器，与 PostgreSQL 数据库的连接配置 URL 通过 config 块定义。

Queryx 还通过在 database 块中更改 adapter 属性和 config 来支持 MySQL 和 SQLite 数据库.

MySQL database 的例子:

```hcl
database "db" {
  adapter = "mysql"

  config "development" {
    url = "root@tcp(localhost:3306)/queryx_test?parseTime=true"
  }
}
```

SQLite database 的例子:

```hcl
database "db" {
  adapter = "sqlite"

  config "test" {
    url = "file:test.sqlite3"
  }
}
```

运行以下命令来自动格式化模式文件:

```sh
queryx format
```

## 数据库管理

运行以下命令来创建 PostgreSQL 数据库。默认情况下，queryx 将读取 development 配置块。可以通过设置 QUERYX_ENV 环境变量来更改其值。

```sh
queryx db:create
```

这与以下命令的作用相同

```sh
QUERYX_ENV=development queryx db:create
```

一旦数据库被创建，queryx 可以自动将数据迁移至`schema.hcl`所定义的模式中:

```sh
queryx db:migrate
```

db:migrate 命令将首先比较数据库的当前状态和 `schema.hcl` 中定义的模式。如果存在差异，它将在 db/migrations 目录中生成 SQL 格式的迁移文件。然后，它将执行生成的迁移文件以使数据库与 `schema.hcl` 中定义的模式保持一致。.

## 代码生成

在样例的 `schema.hcl` 文件中，我们已经定义了生成器

```hcl
generator "client-golang"
```

要生成 Golang 客户端方法，只需执行以下命令：

```sh
queryx generate
```

Queryx 将根据数据库名称在 db 目录中生成 Golang 代码。我们将使用 queryx 生成的代码介绍 CRUD 操作（创建、读取、更新、删除）的基础知识。有关生成的方法的更详细说明，请参阅 ORM API 部分。

首先，我们需要实例化一个客户端对象，它是与数据库交互的入口。

```go
c, err := db.NewClient()
```

Queryx 支持通过变更对象来操作数据库的数据（插入、更新和删除）。对于模式中定义的每个模型，queryx 都会生成一个相应的变更对象，其中包含模型中每个字段的设置方法。这确保了查询的正确性，并使安全地修改数据库中的数据变得容易。

创建一个新的文章:

```go
newPost := c.ChangePost().SetTitle("post title").SetContent("post content")
post, err := c.QueryPost().Create(newPost)
```

Queryx 还支持 Active Record 模式，这使得可以在返回的 queryx 记录上调用 `Update()` 或 `Delete()` 方法。

更新文章:

```go
err := post.Update(c.ChangePost().SetTitle("new post title"))
```

删除文章:

```go
err := post.Delete()
```

Queryx 还支持通过查询进行更新和删除操作.

按标题更新所有文章:

```go
updated, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).UpdateAll(c.ChangePost().SetTitle("new post title"))
```

按标题删除所有文章:

```go
deleted, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).DeleteAll()
```

在 Queryx 中使用主键从数据库中检索数据:

```go
post, err := c.QueryPost().Find(1)
```

按标题检索所有文章:

```go
posts, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).All()
```

从查询结果中获取第一篇文章:

```go
post, err := c.QueryPost().Where(c.PostTitle.EQ("post title")).First()
```

# 定义关联关系

Queryx 支持在模式文件中定义关联关系。它还生成相应的预加载查询方法，以避免 "N+1" 查询问题.

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

## 查询对象

对于模式中定义的每个模型，queryx 都会生成一个相应的查询对象。

```go
q := c.QueryPost()
```

### Finder 方法

查询对象支持以下查找方法:

- Find
- FindBy
- FindBySQL

### Query 方法

查询构建:

- Where
- Limit
- Offset
- Order
- Joins

查询执行:

- All
- First
- Count
- Exists
- UpdateAll
- DeleteAll

### 记录方法

- Update
- Delete

## 事务

Queryx 还支持类型安全的数据库事务，使得可以轻松地安全执行数据库事务。

创建一个事务:

```go
c, err := db.NewClient()
tx := c.Tx()
```

queryx 事务对象的工作方式与 queryx 客户端方法类似，但需要额外的提交调用才能对数据库进行更改。

```go
post, err := tx.QueryPost().Create(tx.ChangPost().SetTitle("post title"))
err := post.Update(tx.ChangePost().SetTitle("new post title"))
if err := tx.Commit() {
  tx.Rollback()
}
```

# 数据类型

queryx 中预定义的数据类型有:

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

# 模式 API

## 命名约定

> **Warning**
> WIP, 请参考以下例子 [here](/internal/integration/postgresql.hcl)

## 时区

默认情况下，每个数据库都使用 "Local" 时区

```hcl
database "db" {
  time_zone = "Local" # this is optional
}
```

指定时区的方法如下:

```hcl
database "db" {
  time_zone = "Africa/Lagos"
}
```

## 环境变量

Queryx 提供了一个方便的功能，可以使用内置的 `env()` 函数从环境变量中读取配置。在生产环境中，从环境变量中读取配置设置是一种常见做法。在下面的示例中，通过将 `QUERYX_ENV` 设置为 production，queryx 将自动从 `DATABASE_URL` 环境变量中读取数据库连接 URL。

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

## 数据库索引

可以通过 `index` 块在模式中声明数据库索引：

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

## 自定义表名

默认情况下，queryx 会以复数形式生成表名。例如，一个名为 User 的模型将有一个名为 users 的表。但是，可以使用模型块中的 table_name 属性自定义此行为。例如:

```hcl
model "User" {
  table_name = "queryx_users"
}
```

在此示例中，queryx 将为 `User` 模型生成表 `queryx_users`。

## 自定义主键

默认情况下，模式中定义的每个模型都会在相应的表中生成一个自增的整数 id 列。可以使用 `primary_key` 块自定义此行为.

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

在这个例子中，Code 模型对应于 codes 表，将具有 type 和 key 作为主键。需要注意的是，自定义主键将影响生成的方法，包括 Find 和 Delete 方法。例如，Code 模型生成的代码中的 Find 方法将不再接受一个整数参数，而是接受两个字符串参数:

```go
func (q *CodeQuery) Find(typ string, key string) (*Code, error)
```

在许多应用程序中，UUID 主键很常见，queryx 支持 UUID 主键的方法如下:

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

默认情况下，queryx CLI 将从当前目录中的 schema.hcl 文件中读取模式。要使用其他模式文件，可以使用 --schema 标志指定文件路径:

```sh
queryx format --schema db.hcl
```

所有可用的命令有:

- `queryx db:create`: 创建数据库
- `queryx db:drop`: 移除数据库
- `queryx db:migrate`: 生成迁移文件并运行待处理的迁移
- `queryx db:migrate:generate`: 生成迁移文件
  <!-- - `queryx db:migrate:status`: list status of each migration -->
  <!-- - `queryx db:rollback`: rollback database -->
  <!-- - `queryx db:version`: print database migration version -->
- `queryx format`: 格式化 hcl 文件
- `queryx generate`: 生成指定客户端代码
- `queryx version`: 显示 queryx 的版本信息

# License

Queryx is licensed under Apache 2.0 as found in the LICENSE file.
