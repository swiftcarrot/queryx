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
    column "emails" {
        type = string
        array = true
    }
  }
}
