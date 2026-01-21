database "db" {
  adapter = "sqlite"

  config "development" {
    url = "sqlite:blog_development.sqlite3"
  }

  generator "client-golang" {}

  model "Post" {
    column "title" {
      type = string
    }
    column "content" {
      type = text
    }
  }
}
