import { ForNow } from "./paths";
// import * as original from "./original";

// UI to get a title goes here...
let title = "this is a new article";

// Slurp the template.
const thetemplate: string = LoadTemplate("entry.md");
console.log(thetemplate);

// Get the date.
const wikidate = Wikidate();
const datepath = ForNow();

// Build the article.

// Make the directory.

// Write the article to the desired location.

// Open the file in iaWriter

console.log(thetemplate);
