import { test, expect, beforeAll } from "vitest";
import { newClient, QXClient } from "./db";

let c: QXClient;

beforeAll(async () => {
  c = newClient();
});

test("string array", async () => {
  let strings = ["test1", "test2"];
  let user = await c.queryUser().create({ strings });
  expect(user.strings).toEqual(strings);
  let row = await c.queryUser().find(user.id);
  expect(row.strings).toEqual(strings);
});

test("text array", async () => {
  let texts = ["test1", "test2"];
  let user = await c.queryUser().create({ texts });
  expect(user.texts).toEqual(texts);
  let row = await c.queryUser().find(user.id);
  expect(row.texts).toEqual(texts);
});

test("integer array", async () => {
  let integers = [1, 2];
  let user = await c.queryUser().create({ integers });
  expect(user.integers).toEqual(integers);
  let row = await c.queryUser().find(user.id);
  expect(row.integers).toEqual(integers);
});
