database "db" {
  adapter = "sqlite"

  config "development" {
    url = "sqlite:test.sqlite3"
  }

  model "User" {
    column "name" {
      type = string
    }
    column "email" {
      type = string
    }
  }
}
