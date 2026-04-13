package summarize

import "sort"

// Result holds a single summarised line and its observed count.
type Result struct {
	Line  string
	Count int
}

// sortResults sorts results by descending count, then ascending line for
// deterministic output when counts are equal.
func sortResults(r []Result) {
	sort.Slice(r, func(i, j int) bool {
		if r[i].Count != r[j].Count {
			return r[i].Count > r[j].Count
		}
		return r[i].Line < r[j].Line
	})
}
