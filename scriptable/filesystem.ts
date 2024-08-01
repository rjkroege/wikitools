// TODO(rjk): insert something here to make sure that the promises are being setup correctly?
export const Markdownext = ".md";

export function JoinPath(
  relpath: string,
  sanetitle: string,
  ext: string,
): string {
  const filepath = [relpath, sanetitle + ext].join("/");
  return filepath;
}

export function AbsPath(
  relpath: string,
  sanetitle: string,
  ext: string,
): string {
  const fm = FileManager.iCloud();
  const abspath = [fm.bookmarkedPath("wiki"), relpath, sanetitle + ext].join(
    "/",
  );
  return abspath;
}

export function WriteFile(abspath: string, payload: string) {
  let fm = FileManager.iCloud();
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
export function LoadTemplate(templatename: string): Promise<string> {
  console.log("loadTemplate: " + templatename);
  let fm = FileManager.iCloud();
  let bpath = fm.bookmarkedPath("wiki");
  let templatepath = bpath + "/templates/" + templatename + Markdownext;

  console.log(fm.listContents(bpath + "/templates"));

  return fm.downloadFileFromiCloud(templatepath).then(() => {
    return fm.readString(templatepath);
  });
}

export function GetTemplateNames(): string[] {
  console.log("GetTemplateNames()");
  const fm = FileManager.iCloud();
  const bpath = fm.bookmarkedPath("wiki");
  const tpath = bpath + "/templates";

  const direntries = fm.listContents(tpath);

  const templatenames = direntries
    .filter((e) => e.endsWith(Markdownext))
    .map((e) => e.substring(0, e.length - Markdownext.length));

  return templatenames;
}
