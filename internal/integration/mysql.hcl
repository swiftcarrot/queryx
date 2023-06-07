database "db" {
  adapter = "mysql"

  config "test" {
    url = "mysql:mysql@tcp(localhost:3306)/queryx_test?parseTime=true"
  }

  generator "client-golang" {}

  model "User" {
    table_name = "queryx_users"
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
      type    = integer
      default = 0
    }

    column "is_admin" {
      type    = boolean
      default = false
    }

    column "payload" {
      type = jsonb
    }

    column "weight" {
      type = float
    }

    column "birth" {
      type = date
      null = true
    }

    column "lunar_birth" {
      type = date
      null = false
    }

    column "login_time" {
      type = datetime
      null = true
    }

    column "lunch_break" {
      type = time
      null = true
    }

  }

  model "Post" {
    has_many "user_posts" {}
    has_many "users" {
      through = "user_posts"
    }

    column "title" {
      type = string
    }
    column "content" {
      type = text
    }
  }

  model "UserPost" {
    table_name = "queryx_user_posts"
    belongs_to "user" {}
    belongs_to "post" {}

    index {
      columns = ["user_id", "post_id"]
      unique  = true
    }
  }

  model "Account" {
    belongs_to "user" {}
    column "name" {
      type = string
    }
    column "id_num" {
      type = integer
    }

  }

  model "Tag" {
    column "name" {
      type = string
    }
  }
}
