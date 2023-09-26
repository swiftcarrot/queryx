# What is Queryx?

::: warning
This project is currently in beta (v0), although it has been battled tested in internal projects. Feel free to [open an issue](https://github.com/swiftcarrot/queryx/issues) or [start a discussion](https://github.com/swiftcarrot/queryx/discussions) if you have any questions.
:::

Queryx is schema-first and type-safe ORM for Go and TypeScript.

- **Schema First**: Queryx automatically migrates the database based on defined models in a queryx schema file.
- **Type Safe**: Queryx generates friendly, type-safe ORM methods and come with autocomplete support and are free from type-related errors.
- **Go and TypeScript**: Queryx supports generating both Go and TypeScript ORM methods.

This project is heavily inspired by [Active Record](https://guides.rubyonrails.org/active_record_basics.html) and [ent](https://entgo.io/). Database schema management is built upon [Atlas](https://atlasgo.io/).
