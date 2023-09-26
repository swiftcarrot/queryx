import { test, expect } from "vitest";
import { newSelect } from "./select";
import { Clause } from "./clause";

test("select", () => {
  let s = newSelect().select("users.*").from("users");
  expect(s.toSQL()).toEqual(["SELECT users.* FROM users", []]);
});

test("select where", () => {
  let s1 = newSelect()
    .select("users.*")
    .from("users")
    .where(new Clause("id = ?", [1]));
  expect(s1.toSQL()).toEqual(["SELECT users.* FROM users WHERE id = ?", [1]]);

  s1.where(new Clause("name = ?", ["test"]));
  expect(s1.toSQL()).toEqual([
    "SELECT users.* FROM users WHERE (id = ?) AND (name = ?)",
    [1, "test"],
  ]);
});
