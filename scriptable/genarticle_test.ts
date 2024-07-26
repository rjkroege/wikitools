// [bun test – Test runner | Bun Docs](https://bun.sh/docs/cli/test)
import { expect, test } from "bun:test";
import { GenerateArticle } from "./genarticle";

test("genarticles", () => {
  expect(GenerateArticle("Write something", "My Title", "Date")).toBe(
    "---\ntitle: My Title\ndate: Date\ntags: \n---\n\nWrite something\n",
  );
});
