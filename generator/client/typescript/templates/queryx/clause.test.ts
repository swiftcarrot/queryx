// Code generated by queryx, DO NOT EDIT.

import { test, expect } from "vitest";
import { Clause } from "./clause";

test("clauseAndOr", () => {
  let c1 = new Clause("a = ?", [1]);
  let c2 = new Clause("b = ?", ["x"]);

  let c3 = c1.and(c2);
  expect(c3.fragment).toEqual("(a = ?) AND (b = ?)");
  expect(c3.args).toEqual([1, "x"]);

  let c4 = c1.or(c2);
  expect(c4.fragment).toEqual("(a = ?) OR (b = ?)");
  expect(c4.args).toEqual([1, "x"]);

  let c5 = c1.and(c1.or(c2));
  expect(c5.fragment).toEqual("(a = ?) AND ((a = ?) OR (b = ?))");
  expect(c5.args).toEqual([1, 1, "x"]);

  expect(c1.fragment).toEqual("a = ?");
  expect(c1.args).toEqual([1]);
  expect(c2.fragment).toEqual("b = ?");
  expect(c2.args).toEqual(["x"]);
});
