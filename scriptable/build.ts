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

// Define the targets and their entry points
const targets = [
  {
    entrypoint: "./index.ts",
    outputName: "wikitools",
  },
  {
    entrypoint: "./websummary.ts", // Replace with your second entry point
    outputName: "websummary",
  },
];

// Build each target
for (const target of targets) {
  const thebuild = await Bun.build({
    entrypoints: [target.entrypoint],
    naming: `[dir]/${target.outputName}.[ext]`,
  });

  for (const output of thebuild.outputs) {
    const blob = await output;
    await mkdir(thepath, { recursive: true });
    const thefile = path.join(thepath, output.path);

    // Attempt to unlink the file, ignore error if it doesn't exist
    try {
      await unlink(thefile);
    } catch (error) {
      if (error.code !== 'ENOENT') {
        // Log the error if it's not a "file not found" error
        console.error(`Error unlinking file: ${thefile}`, error);
      }
    }

    const fd = Bun.file(thefile);
    const fdw = fd.writer();
    fdw.write(contents);
    fdw.write(await blob.arrayBuffer());
    fdw.end();
  }
}
