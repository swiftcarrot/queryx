import { test, expect, beforeAll } from "vitest";
import { newClientWithEnv, QXClient } from "./db";

let c: QXClient;

beforeAll(async () => {
  c = await newClientWithEnv("test");
});

test("queryOne", async () => {
  let user = await c.queryUser().create({
    name: "test",
  });
  let row = await c.queryOne(
    "select id as user_id from users where id = ?",
    user.id
  );
  expect(row.user_id).toEqual(user.id);
});

test("query", async () => {
  let user1 = await c.queryUser().create({ name: "test1" });
  let user2 = await c.queryUser().create({ name: "test2" });
  let rows = await c.query(
    "select name as user_name from users where id in (?)",
    [user1.id, user2.id]
  );
  expect(rows).toEqual([{ user_name: "test1" }, { user_name: "test2" }]);
});

test("exec", async () => {
  let user = await c.queryUser().create({ name: "test" });
  let updated = await c.exec(
    "update users set name = ? where id = ?",
    "test1",
    user.id
  );
  expect(updated).toEqual(1);
  let deleted = await c.exec("delete from users where id = ?", user.id);
  expect(deleted).toEqual(1);
});

test("create", async () => {
  let user = await c.queryUser().create({ name: "test", type: "admin" });
  expect(user.name).toEqual("test");
  expect(user.type).toEqual("admin");
  expect(user.id).toBeGreaterThan(0);
});

test("transaction", async () => {
  let tag1 = await c.queryTag().create({ name: "tag1" });
  expect(tag1.name).toEqual("tag1");

  let total1 = await c.queryTag().count();
  let tx = await c.tx();

  tag1 = await tx.queryTag().find(tag1.id);
  await tag1.update({ name: "tag1-updated" });

  await tx.queryTag().create({ name: "tag2" });
  await tx.queryTag().create({ name: "tag3" });

  let total2 = await c.queryTag().count();
  expect(total2).toEqual(total1);

  let total3 = await c.queryTag().count();
  expect(total3).toEqual(total1 + 2);

  tag1 = c.queryTag().find(tag1.id);
  expect(tag1.name).toEqual("tag1");

  await tx.commit();

  let total4 = c.queryTag().count();
  expect(total4).toEqual(total1 + 2);

  tag1 = c.queryTag().find(tag1.id);
  expect(tag1.name).toEqual("tag1-updated");
});
