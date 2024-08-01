import { ForNow, SaneFileName, Wikidate } from "./paths";
import {
  WriteFile,
  Makedir,
  LoadTemplate,
  Markdownext,
  JoinPath,
  AbsPath,
  GetTemplateNames,
} from "./filesystem";
import { GenerateArticle } from "./genarticle";
import { OpenInEditor } from "./iawriter";
import { ShowNewArticleDialog } from "./alert";
// Need a special statement to import types.
import type { WikiArticleParms } from "./alert";

const alltemplatenames: string[] = await GetTemplateNames();
console.log(alltemplatenames);

const parms: WikiArticleParms = await ShowNewArticleDialog(alltemplatenames);

const title = parms.title;

// Slurp the template.
const thetemplate: string = await LoadTemplate(parms.template);

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
