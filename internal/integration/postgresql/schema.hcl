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
    column "strings" {
      type  = string
      array = true
    }
    column "integers" {
      type  = integer
      array = true
    }
    column "texts" {
      type  = text
      array = true
    }
  }
}
