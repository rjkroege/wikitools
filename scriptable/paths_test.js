// [bun test – Test runner | Bun Docs](https://bun.sh/docs/cli/test)
import { expect, test } from "bun:test";
import { forNowImpl } from "./paths.js";

test("1995-12-17T03:24:00", () => {
  expect(forNowImpl(new Date("1995-12-17T03:24:00"))).toBe("1995/12-Dec/17");
});

test("1995-03-17T03:24:00", () => {
  expect(forNowImpl(new Date("1995-03-17T03:24:00"))).toBe("1995/03-Mar/17");
});
