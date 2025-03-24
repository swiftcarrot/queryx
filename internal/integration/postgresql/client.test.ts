import { test, expect, beforeAll } from "vitest";
import { newClient, QXClient } from "./db";

let c: QXClient;

beforeAll(async () => {
  c = newClient();
});

test("string array", async () => {
  let emails = ["test1@example.com", "test2@example.com"];
  let user = await c.queryUser().create({ emails: emails });
  expect(user.emails).toEqual(emails);
  let row = await c.queryUser().find(user.id);
  expect(row.emails).toEqual(emails);
});
