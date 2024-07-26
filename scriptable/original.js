import { ForNow } from "./paths";

// Make a call to the app in which to write journal articles.
function openIaWriter(title, contents) {
  let iawurl = new CallbackURL("iawriter://new");
  // 	compute the right path
  iawurl.addParameter("path", title);
  iawurl.addParameter("text", contents);
  iawurl.open();
}

// templates are stored with respect to the wiki bookmark.
// path is the name of the template in the templates directory.
// returns the template.
function LoadTemplate(path) {
  console.log("loadTemplate: " + path);
  let fm = FileManager.iCloud();
  let bpath = fm.bookmarkedPath("wiki");
  let templatepath = bpath + "/templates/" + path + ".md";

  if (fm.fileExists(bpath)) {
    console.log(bpath + " exists");
  }
  if (fm.fileExists(bpath + "/templates")) {
    console.log(bpath + "/templates exists");
  }
  if (fm.fileExists(templatepath)) {
    console.log(templatepath + " exists");
  }
  console.log(fm.listContents(bpath + "/templates"));
  let s = fm.readString(templatepath);
  console.log(s);

  let fp = fm.downloadFileFromiCloud(templatepath);
  fp.then(() => {
    if (fm.fileExists(templatepath)) {
      console.log(templatepath + " exists");
    } else {
      console.log(templatepath + " doesn't exist");
    }
    console.log("hi");
    let templatestring = fm.readString(templatepath);
    console.log("load: " + templatestring);
  });

  return templatestring;
}

//
function writeFileToSomewhere(targetpath) {
  let fm = FileManager.iCloud();
  let bpath = fm.bookmarkedPath("wiki");
}

// pickTemplate();

// listDirectory();

// /Locations/wiki/unsorted/

// openIaWriter(kNewJournalArticlePath + "test1.md" , foo);

// i don't really want this. i instead want to do load tbe right template.

// idea: ask for the file to open. save config.

// let s = loadTemplate("journal_pm");
// let timer = Timer.schedule(10*1000, false, ()=>{ console.log("end timer")});
