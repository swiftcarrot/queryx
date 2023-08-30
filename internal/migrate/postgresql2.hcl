database "db" {
  adapter   = "postgresql"
  time_zone = "Asia/Shanghai"

  config "development" {
    url = "postgresql://postgres:postgres@localhost:5432/queryx_test?sslmode=disable"
  }
  model "User" {
    timestamps = false

    column "name" {
      type = string
    }
  }
}