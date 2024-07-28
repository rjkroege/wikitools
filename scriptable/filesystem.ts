// TODO(rjk): insert something here to make sure that the promises are being setup correctly?
const markdownext = ".md";

export function WriteFile(relpath: string, sanetitle: string, payload: string) {
  let fm = FileManager.iCloud();
  const abspath = [fm.bookmarkedPath("wiki"), relpath, sanetitle + markdownext].join("/");
  fm.writeString(abspath, payload);
}

export function Makedir(relpath: string) {
  const fm = FileManager.iCloud();
  const bpath = fm.bookmarkedPath("wiki");
  const abspath = bpath + "/" + relpath;
  fm.createDirectory(abspath, true);
}


// this is getting transpiled in a way that I don't comprehend.
// templates are stored with respect to the wiki bookmark.
// TODO(rjk): needs to return a promse
// This code is dumb. Fix it
export function LoadTemplate(templatename: string) : Promise<string> {
  console.log("loadTemplate: " + templatename);
  let fm = FileManager.iCloud();
  let bpath = fm.bookmarkedPath("wiki");
  let templatepath = bpath + "/templates/" + templatename + markdownext;

  console.log(fm.listContents(bpath + "/templates"));

  return fm.downloadFileFromiCloud(templatepath).then(
	() => {
		return fm.readString(templatepath);
	});
   
}

