database "db" {
  adapter   = "mysql"
  time_zone = "Asia/Shanghai"

  config "development" {
    url = "mysql://root:@127.0.0.1:3306/test"
  }
  config "test" {
    url = env("DATABASE_URL")
  }

  generator "client-typescript" {
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
