## Database Index

Database index can be declared in schema via the `index` block:

```hcl
model "UserPost" {
  belongs_to "user" {}
  belongs_to "post" {}

  index {
    columns = ["user_id", "post_id"]
    unique  = true
  }
}
```
