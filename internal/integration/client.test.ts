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

test("time", async () => {});

test("date", async () => {});

test("datetime", async () => {});

test("timestamps", async () => {
  let user = await c.queryUser().create({});
  expect(user.createdAt).not.toBeNull();
  expect(user.updatedAt).not.toBeNull();
  expect(user.createdAt).toEqual(user.updatedAt);

  await user.update({ name: "new name" });
  expect(user.updatedAt).toBeGreaterThan(user.createdAt);
});

test("uuid", async () => {});

test("null", async () => {
  let user = await c.queryUser().create({ name: null });
  expect(user.name).toBeNull();

  user = await c.queryUser().find(user.id);
  expect(user.name).toBeNull();
});

test("json", async () => {});

test("primaryKey", async () => {
  await c.queryCode().deleteAll();
  let code = await c.queryCode().create({ type: "type", key: "key" });
  expect(code.type).toEqual("type");
  expect(code.key).toEqual("key");

  expect(async () => {
    await c.queryCode().create({ type: "type", key: "key" });
  }).toThrowError();

  code = await c.queryCode().find("type", "key");
  expect(code.type).toEqual("type");
  expect(code.key).toEqual("key");

  c.queryClient().deleteAll();
  let client = await c.queryClient().create({ name: "client" });
  expect(client.name).toEqual("client");

  let deleted = await c.queryClient().delete("client");
  expect(deleted).toEqual(1);
});

test("boolean", async () => {
  let user = await c.queryUser().create({ isAdmin: true });
  expect(user.isAdmin).toBe(true);
});

test("exists", async () => {
  await c.queryClient().deleteAll();
  let exists = await c.queryClient().exists();
  expect(exists).toBe(false);

  await c.queryClient().create({ name: "client" });
  exists = await c.queryClient().exists();
  expect(exists).toBe(true);
});

test("belongsTo", async () => {
  let author = await c.queryUser().create({ name: "author" });
  let post = await c
    .queryPost()
    .create({ title: "post title", authorID: author.id });
  post = await c.queryPost().preloadAuthor().find(post.id);
  expect(post.author.id).toEqual(author.id);
});

test("allEmpty", async () => {
  await c.queryUser().deleteAll();
  let users = await c.queryUser().all();
  expect(users).toEqual([]);
});

test("inEmpty", async () => {
  await c.queryUser().deleteAll();
  let users = await c.queryUser().where(c.userID.in([])).all();
  expect(users).toEqual([]);

  users = await c
    .queryUser()
    .where(c.userID.in([]).and(c.userID.in([])).and(c.userID.in([])))
    .all();
  expect(users).toEqual([]);
});

test("hasManyEmpty", async () => {
  let user = await c.queryUser().create({ name: "user" });
  expect(user.userPosts).toBeNull();
  expect(user.posts).toBeNull();

  user = await c.queryUser().preloadUserPosts().find(user.id);
  expect(user.userPosts).toEqual([]);

  user = await c.queryUser().preloadPosts().find(user.id);
  expect(user.posts).toEqual([]);
  expect(user.userPosts).toEqual([]);
});

test("hasOne", async () => {
  let user = await c.queryUser().create({ name: "user" });
  let account = await c
    .queryAccount()
    .create({ name: "account", userID: user.id });
  user = await c.queryUser().preloadAccount().find(user.id);
  expect(user.account.name).toEqual(account.name);
});

test("preload", async () => {
  let user1 = await c.queryUser().create({ name: "user1" });
  let post1 = await c.queryPost().create({ title: "post1" });
  let post2 = await c.queryPost().create({ title: "post2" });
  let account1 = await c
    .queryAccount()
    .create({ name: "account1", userID: user1.id });
  let userPost1 = await c
    .queryUserPost()
    .create({ userID: user1.id, postID: post1.id });
  let userPost2 = await c
    .queryUserPost()
    .create({ userID: user1.id, postID: post2.id });
  let user = await c.queryUser().preloadPosts().preloadAccount().find(user1.id);
  expect(user.account.id).toEqual(account1.id);
  expect(user.userPosts.length).toEqual(2);
  expect(user.userPosts[0].id).toEqual(userPost1.id);
  expect(user.userPosts[1].id).toEqual(userPost2.id);

  expect(user.posts.length).toEqual(2);
  expect(user.posts[0].id).toEqual(post1.id);
  expect(user.posts[1].id).toEqual(post2.id);

  let post = await c.queryPost().preloadUserPosts().find(post1.id);
  expect(post.userPosts.length).toEqual(1);
  expect(post.userPosts[0].id).toEqual(userPost1.id);
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

  let total3 = await tx.queryTag().count();
  expect(total3).toEqual(total1 + 2n);

  tag1 = await c.queryTag().find(tag1.id);
  expect(tag1.name).toEqual("tag1");

  await tx.commit();

  let total4 = await c.queryTag().count();
  expect(total4).toEqual(total1 + 2n);

  tag1 = await c.queryTag().find(tag1.id);
  expect(tag1.name).toEqual("tag1-updated");
});

test("changeJSON", async () => {});

test("modelString", async () => {});
