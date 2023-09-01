database "db" {
  adapter   = "postgresql"

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
