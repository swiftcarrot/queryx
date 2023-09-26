# Custom Time Zone

By default, each database uses the "Local" time zone.

```hcl{2}
database "db" {
  time_zone = "Local" # this is optional
}
```

To specify time zone:

```hcl{2}
database "db" {
  time_zone = "Africa/Lagos"
}
```
