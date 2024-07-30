// Make a call to the app in which to write journal articles.

export function OpenInEditor(librarylocation: string, relpath: string) {
  // Implication of [URL Commands – iA](https://ia.net/writer/support/help/url-commands) is that
  // the path component works like a wikilink.
  let iawurl = new CallbackURL("ia-writer://open");
  // Compute the right path
  const path = `${librarylocation}: ${relpath}`;
  log(`iawurl: ${path}`);
  iawurl.addParameter("path", path);
  iawurl.addParameter("edit", "true");
  iawurl.open();
}
