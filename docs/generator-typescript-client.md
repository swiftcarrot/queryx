# generator-typescript-client

Queryx supports generating ORM methods in TypeScript for your schema. To enable this generator, add the following code to your schema file:

```hcl
database "db" {
  generator "typescript-client" {}
}
```

Before using the generated methods, you need to install some external dependencies. Queryx relies on `date-fns` for handling dates in TypeScript. Additionally, depending on your database, you will need to install specific packages. Here are the installation commands for different databases:

For PostgreSQL:

```sh
npm install pg @types/pg
```

For MySQL:

```sh
npm install mysql2 @types/node
```

For SQLite:

```sh
npm install better-sqlite3
```

Once you have installed the required dependencies, you can start using the generated methods. Unlike the generated Golang code, the TypeScript version does not require a Change object for modifying database records. Here's an example of how to use the generated methods:

```typescript
import { newClient } from "./db";

let client = await newClient();
let user = await client.queryUser().create({ name: "user name" });
await user.update({ name: "new user name" });
```

In the above example, `newClient()` creates a new client object. You can then use the generated `queryUser()` method to perform various operations on the `User` table, such as creating a new user or updating an existing user.
