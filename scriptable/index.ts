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
import { OpenInEditor, NewInEditor } from "./iawriter";
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

const relpath = JoinPath(datepath, SaneFileName(title), Markdownext);
const librarylocation = "wiki";
const iapath = `${librarylocation}: ${relpath}`;

NewInEditor(iapath, articletext);
