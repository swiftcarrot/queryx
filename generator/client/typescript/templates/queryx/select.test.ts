import { test, expect } from "vitest";
import { newSelect } from "./select";

test("select", () => {
  let s = newSelect().select("users.*").from("users");
  expect(s.toSQL()).toEqual(["select users.* from users", []]);
});
