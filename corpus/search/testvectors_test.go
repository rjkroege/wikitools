package search

const pathoffset = 5

var wikitests = []nameIndexTestVector{
	{input_location: "wiki", input_wikitextfile: "Saturday.md", want: []string{
		"wiki/2023/02-Feb/28/Saturday.md",
		"wiki/2023/05-May/6/Saturday.md",
		"wiki/unsorted/Saturday.md",
	}, want_err: nil},
	{input_location: "wiki", input_wikitextfile: "Coffee.md", want: []string{"wiki/2023/12-Dec/Coffee.md"}, want_err: nil},
	{input_location: "wiki", input_wikitextfile: "Thursday morning.md", want: []string{"wiki/unsorted/Thursday morning.md"}, want_err: nil},
	{input_location: "wiki", input_wikitextfile: "missing file name.md", want: []string{}, want_err: nil},
	{input_location: "wiki/2023", input_wikitextfile: "Saturday.md", want: []string{
		"wiki/2023/02-Feb/28/Saturday.md",
		"wiki/2023/05-May/6/Saturday.md",
	}, want_err: nil},
	{input_location: "wiki/2023/02-Feb", input_wikitextfile: "Saturday.md", want: []string{
		"wiki/2023/02-Feb/28/Saturday.md",
	}, want_err: nil},
}
