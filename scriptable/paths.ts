// Export ForNow as a function visible outside of this class.
import date from "date-and-time";

// Get the relative path for today's date.
export function ForNow() {
  let now = new Date();
  return forNowImpl(now);
}

// TODO(rjk): This code is dumb about timezones!
export function forNowImpl(now: Date): string {
  return date.format(now, "YYYY/MM-MMM/DD");
}

const dewhiter = /\s/gu;
// Regexp derived from https://github.com/sindresorhus/filename-reserved-regex/tree/main, MIT license
const cleaner = /[<>:"/\\|?*\u0000-\u001F^]/gu;

// Makes name into a sane name.
// TODO(rjk): Add the extension.
export function SaneFileName(filename: string): string {
  const saner = filename.replace(cleaner, "_");
  return saner.replace(dewhiter, "-");
}

export function Wikidate(): string {
  const now = new Date();
  return wikidateimpl(now);
}

export function wikidateimpl(now: Date): string {
  // unixlikezoned  = "Mon _2 Jan 2006, 15:04:05 -0700"
  return date.format(now, "ddd DD MMM YYYY, HH:mm:ss Z");
}
