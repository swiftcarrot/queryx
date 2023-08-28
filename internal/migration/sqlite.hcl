database "db" {
  adapter = "sqlite"

  config "test" {
    url = "sqlite:test.sqlite3"
  }

  model "User" {
    column "name" {
      type = string
    }
  }
}
