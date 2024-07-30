// Build script

// Bun lets me import text files and they become variables.
// NB: This construct however appears to confuse the typescript compiler.
// TODO(rjk): Fix this.
import contents from "./header.txt";

// This imports the entire package.
import path from "node:path";

// This imports a single function I think.
import { mkdir, unlink } from "node:fs/promises";

// Logs something.
const home: string = process.env.HOME ?? "/home/me";
const thepath = path.join(
  home,
  "Library/Mobile Documents/iCloud~dk~simonbs~Scriptable/Documents",
);
console.log("thepath", thepath);

const thebuild = await Bun.build({
  entrypoints: ["./index.ts"],
  // Don't set outdir so that Bun writes nothing itself.
  // But keep the naming.
  naming: "[dir]/wikitools.[ext]",
});

for (const output of thebuild.outputs) {
  const blob = await output;
  await mkdir(thepath, { recursive: true });
  const thefile = path.join(thepath, output.path);
  await unlink(thefile);
  const fd = Bun.file(thefile);
  const fdw = fd.writer();
  fdw.write(contents);
  fdw.write(await blob.arrayBuffer());
  fdw.end();
}
