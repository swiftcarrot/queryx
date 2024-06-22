database "db" {
  adapter   = "sqlite"
  time_zone = "Asia/Shanghai"

  config "development" {
    url = "sqlite:test.sqlite3"
  }
  config "test" {
    url = env("DATABASE_URL")
  }

  generator "client-golang" {
    test = true
  }

  model "Model" {
    timestamps = false
    column "name" {
      type = string
    }
    column "title" {
      type = string
    }
    column "fax" {
      type = string
    }
    column "web" {
      type = string
    }
    column "age" {
      type = bigint
    }
    column "righ" {
      type = boolean
    }
    column "counter" {
      type = integer
    }
  }

}
