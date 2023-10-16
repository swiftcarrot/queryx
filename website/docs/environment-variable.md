# Environment Variable

Queryx provides a convenient feature for reading from environment variables using the built-in `env()` HCL function. It is a common practice for applications to read configuration settings from environment variables in production environments. In the following example, by setting `QUERYX_ENV` to `production`, queryx will automatically read the database connection URL from the `DATABASE_URL` environment variable.

```hcl{6-8}
database "db" {
  config "development" {
    url = "postgres://postgres:postgres@localhost:5432/blog_development?sslmode=disable"
  }

  config "production" {
    url = env("DATABASE_URL")
  }
}
```

