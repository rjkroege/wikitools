// Export ForNow as a function visible outside of this class.

const months : string[]  = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec"
];

// Get the relative path for today's date.
export function ForNow() {
  let now = new Date();
  return forNowImpl(now);
}

export function forNowImpl(now: Date) : string {
  let year = now.getFullYear();
  let month0 = now.getMonth();
  let month = month0 + 1;
  let day = now.getDate();

  // zero-pad.
  let finalmonth =
    (month < 10 ? "0" + month : "" + month) + "-" + months[month0];

  // is there a joining...
  return [year, finalmonth, day].join("/");
}

const dewhiter = /\s/ug;
// Regexp derived from https://github.com/sindresorhus/filename-reserved-regex/tree/main, MIT license
const cleaner = /[<>:"/\\|?*\u0000-\u001F^]/ug;

// Makes name into a sane name.
export function SaneFileName(filename: string) : string {
	const saner = filename.replace(cleaner, "_");
	return saner.replace(dewhiter, "-");
}
