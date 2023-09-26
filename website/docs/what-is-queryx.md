# What is Queryx?

::: warning
This project is currently in beta (v0), although it has been battled tested in internal projects. Feel free to [open an issue](https://github.com/swiftcarrot/queryx/issues) or [start a discussion](https://github.com/swiftcarrot/queryx/discussions) if you have any questions.
:::

Queryx is schema-first and type-safe ORM for Go and TypeScript.

- **Schema First**: Define application models in a queryx schema file, and it can automatically synchronize with the database structure.
- **Type Safe**: Queryx generates friendly and type-safe ORM methods based on the schema, which come with autocomplete support and are free from type-related errors.
- **Go and TypeScript**: Queryx can generates both Go and TypeScript ORM methods based on schema.

This project is heavily inspired by [Active Record](https://guides.rubyonrails.org/active_record_basics.html) and [ent](https://entgo.io/). Database schema management is built upon [Atlas](https://atlasgo.io/).
