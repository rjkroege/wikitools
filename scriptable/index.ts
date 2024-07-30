import { ForNow, SaneFileName, Wikidate } from "./paths";
import {
  WriteFile,
  Makedir,
  LoadTemplate,
  Markdownext,
  JoinPath,
  AbsPath,
} from "./filesystem";
import { GenerateArticle } from "./genarticle";
import { OpenInEditor } from "./iawriter";

// UI to get a title goes here...
let title = "this is a new article";

// Slurp the template.
const thetemplate: string = await LoadTemplate("entry");
console.log(thetemplate);

// Get the date.
const wikidate = Wikidate();
const datepath = ForNow();

// Build the article.
const articletext = GenerateArticle(thetemplate, title, wikidate);

// Make the directory.
Makedir(datepath);

// TODO(rjk): Perhaps further refactor the pathing code.
// Write the article to the desired location.
const abspath = AbsPath(datepath, SaneFileName(title), Markdownext);
WriteFile(abspath, articletext);

OpenInEditor("wiki", JoinPath(datepath, SaneFileName(title), Markdownext));

console.log(thetemplate);
