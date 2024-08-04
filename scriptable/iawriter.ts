// Make a call to the app in which to write journal articles.

export function OpenInEditor(iapath: string) {
  // Implication of [URL Commands – iA](https://ia.net/writer/support/help/url-commands) is that
  // the path component works like a wikilink.
  let iawurl = new CallbackURL("ia-writer://open");
  iawurl.addParameter("path", iapath);
  iawurl.addParameter("edit", "true");
  iawurl.open();
}

export function NewInEditor(iapath: string, articletext: string) {
  let iawurl = new CallbackURL("ia-writer://new");
  iawurl.addParameter("path", iapath);
  iawurl.addParameter("text", articletext);
  iawurl.addParameter("edit", "true");
  iawurl.open();
}
