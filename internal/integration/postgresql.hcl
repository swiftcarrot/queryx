database "db" {
  adapter   = "postgresql"
  time_zone = "Asia/Shanghai"

  config "development" {
    url = "postgresql://postgres:postgres@localhost:5432/queryx_test?sslmode=disable"
  }
  config "test" {
    url = env("DATABASE_URL")
  }

  generator "client-golang" {
    test = true
  }
  generator "client-typescript" {
    test = true
  }

  model "User" {
    has_one "account" {}
    has_many "user_posts" {}
    has_many "posts" {
      through = "user_posts"
    }

    column "name" {
      type = string
    }
    column "type" {
      type = string
    }
    column "email" {
      type = string
    }
    column "age" {
      type = integer
    }
    column "is_admin" {
      type = boolean
    }
    column "payload" {
      type = jsonb
    }
    column "weight" {
      type = float
    }
    column "date" {
      type = date
    }
    column "datetime" {
      type = datetime
    }
    column "time" {
      type = time
    }
    column "uuid" {
      type = uuid
    }
  }

  model "Post" {
    has_many "user_posts" {}
    has_many "users" {
      through = "user_posts"
    }
    belongs_to "author" {
      model_name = "User"
    }

    column "title" {
      type = string
    }
    column "content" {
      type = text
    }
    column "payload" {
      type = json
    }
    column "tags" {
      type  = string
      array = true
    }
  }

  model "UserPost" {
    belongs_to "user" {}
    belongs_to "post" {}

    index {
      columns = ["user_id", "post_id"]
      unique  = true
    }
  }

  model "Account" {
    belongs_to "user" {
      index = true
      null  = false
    }

    column "name" {
      type = string
    }
    column "id_num" {
      type = integer
    }
  }

  model "Tag" {
    timestamps = false

    column "name" {
      type = string
    }

    index {
      columns = ["name"]
      unique  = true
    }
  }

  model "Code" {
    table_name          = "queryx_codes"
    default_primary_key = false
    timestamps          = false

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

  model "Client" {
    default_primary_key = false
    timestamps          = false

    column "name" {
      type = string
    }
    column "float" {
      type = float
    }
  }

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
}
