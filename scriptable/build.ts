// Build script

// As of July 2024, bun, node, TypeScript ecosystem is new to me. So I'm
// adding extra comments to capture
import contents from "./header.txt";

// This imports the entire package.
import path from "node:path";

// This imports a single function I think.
import { mkdir } from "node:fs/promises";

// Logs something.
console.log("building wikitools package");
const thepath = "build";

const thebuild = await Bun.build({
  entrypoints: ["./index.ts"],
  // Don't set outdir so that Bun writes nothing itself.
  // But keep the naming.
  naming: "[dir]/wikitools.[ext]",
});

for (const output of thebuild.outputs) {
  const blob = await output;
  await mkdir(thepath, { recursive: true });
  const fd = Bun.file(path.join(thepath, output.path));
  const fdw = fd.writer();
  fdw.write(contents);
  fdw.write(await blob.arrayBuffer());
  fdw.end();
}
