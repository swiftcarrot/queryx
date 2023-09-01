database "db" {
  adapter   = "mysql"

  config "development" {
    url = "mysql://root:@127.0.0.1:3306/queryx_test"
  }

  model "User" {
    timestamps = false

    column "name" {
      type = string
    }
    column "email" {
      type = string
    }
  }
}
