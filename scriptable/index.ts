import { ForNow, SaneFileName, Wikidate} from "./paths";
import { WriteFile, Makedir, LoadTemplate } from "./filesystem";
import { GenerateArticle } from "./genarticle";

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

// Write the article to the desired location.
WriteFile(datepath, SaneFileName(title), articletext);

console.log(thetemplate);
