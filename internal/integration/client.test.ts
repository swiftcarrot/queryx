import { test, expect, beforeAll } from "vitest";
import { newClientWithEnv, QXClient, UserChange } from "./db";

let c: QXClient;

beforeAll(async () => {
  c = newClientWithEnv("test");
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

test("find", async () => {
  let tag = await c.queryTag().create({ name: "test" });
  tag = await c.queryTag().find(tag.id);
  expect(tag.name).toEqual("test");

  await expect(async () => {
    await c.queryTag().find(Number(tag.id) + 1);
  }).rejects.toThrowError("Record not found");
});

test("order", async () => {
  await c.queryTag().deleteAll();
  let tag1 = await c.queryTag().create({ name: "tag1" });
  let tag2 = await c.queryTag().create({ name: "tag2" });

  let tags = await c.queryTag().order("id desc").all();
  expect(tags).toEqual([tag2, tag1]);
});

test("time", async () => {
  let user = await c.queryUser().create({ time: "12:10:09" });
  expect(user.time).toEqual("12:10:09");
});

test("date", async () => {
  let user = await c.queryUser().create({ date: "2012-11-10" });
  expect(user.date).toEqual("2012-11-10");
});

test("datetime", async () => {
  let s1 = "2012-11-10 09:08:07";
  let user = await c.queryUser().create({ datetime: s1 });
  expect(user.datetime).toEqual(s1);

  user = await c
    .queryUser()
    .where(
      c.userID.eq(user.id).and(c.userDatetime.ge(s1)).and(c.userDatetime.le(s1))
    )
    .first();
  expect(user.datetime).toEqual(s1);

  let s2 = "2012-11-10 09:08:07.654";
  user = c.queryUser().create({ datetime: s2 });
  expect(user.datetime).toEqual(s2);
});

test("timestamps", async () => {
  let user = await c.queryUser().create({});
  expect(user.createdAt).not.toBeNull();
  expect(user.updatedAt).not.toBeNull();
  expect(user.createdAt).toEqual(user.updatedAt);

  await user.update({ name: "new name" });
  expect(user.updatedAt).toBeGreaterThan(user.createdAt);
});

test("uuid", async () => {
  let user = await c.queryUser().create({});
  expect(user.uuid).toBeNull();

  let uuid1 = "c7e5b9af-0499-4eca-a7e6-77e10d56987b";
  await user.update({ uuid: uuid1 });
  expect(user.uuid).toEqual(uuid1);

  let uuid2 = "a81e44c5-7e18-4dfe-b9b3-d9280629d2ef";
  let device = await c.queryDevice().create({ id: uuid2 });
  expect(device.id).toEqual(uuid2);

  device = await c.queryDevice().find(uuid2);
  expect(device.id).toEqual(uuid2);
});

test("null", async () => {
  let user = await c.queryUser().create({ name: null });
  expect(user.name).toBeNull();

  user = await c.queryUser().find(user.id);
  expect(user.name).toBeNull();
});

test("json", async () => {
  let payload = {
    theme: "dark",
    height: 170,
    weight: 65,
  };
  let user = await c.queryUser().create({ payload });
  expect(user.payload.theme).toEqual(payload.theme);
  expect(user.payload.height).toEqual(payload.height);
  expect(user.payload.weight).toEqual(payload.weight);
});

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

test("changeJSON", async () => {
  let userChange = new UserChange({
    name: "user name",
    isAdmin: false,
  });
  expect(userChange.name).toEqual("user name");
  expect(userChange.isAdmin).toBe(false);
});

test("modelJSON", async () => {
  let tag = await c.queryTag().create({ name: "test" });
  expect(JSON.stringify(tag)).toEqual(`{"id":${tag.id},"name":"test"}`);
});

test("modelString", async () => {
  await c.queryCode().deleteAll();
  let code = await c.queryCode().create({
    type: "code type",
    key: "code key",
  });
  expect(code.toString()).toEqual(
    `(Code type: "${code.type}", key: "${code.key}")`
  );
});
