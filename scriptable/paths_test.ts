// [bun test – Test runner | Bun Docs](https://bun.sh/docs/cli/test)
import { expect, test } from "bun:test";
import { forNowImpl, SaneFileName, wikidateimpl } from "./paths";

test("1995-12-17T03:24:00", () => {
  expect(forNowImpl(new Date("1995-12-17T03:24:00"))).toBe("1995/12-Dec/17");
});

test("1995-03-17T03:24:00", () => {
  expect(forNowImpl(new Date("1995-03-17T03:24:00"))).toBe("1995/03-Mar/17");
});

test("SaneFileName: hello", () => {
  expect(SaneFileName("hello")).toBe("hello");
});

test("SaneFileName: hello world", () => {
  expect(SaneFileName("hello world")).toBe("hello-world");
});

test("SaneFileName: hello world/", () => {
  expect(SaneFileName("hello world/")).toBe("hello-world_");
});

test("SaneFileName: ^*郵便局/<<", () => {
  expect(SaneFileName("^*郵便局/<<")).toBe("__郵便局___");
});

test("wikidate 1995-12-17T03:24:00", () => {
  expect(wikidateimpl(new Date("1995-12-17T03:24:00"))).toBe(
    "Sun 17 Dec 1995, 03:24:00 UTC",
  );
});

test("wikidate 1995-03-04T03:24:00", () => {
  expect(wikidateimpl(new Date("1995-03-04T11:24:15+04:00"))).toBe(
    "Sat 04 Mar 1995, 07:24:15 UTC",
  );
});
