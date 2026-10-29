// Package wordcount provides a bounded-line, Unicode whitespace text counter.
package wordcount

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"slices"
	"strings"
)

// Frequency is one lower-cased whitespace-delimited token and its occurrence count.
type Frequency struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// Stats describes text. Lines follows bufio.ScanLines, including a final nonempty line.
type Stats struct {
	Lines       int         `json:"lines"`
	Words       int         `json:"words"`
	Unique      int         `json:"unique"`
	Frequencies []Frequency `json:"frequencies"`
}

// Count reads text with a maximum scanner token size of 1 MiB.
// Tokens are separated by Unicode whitespace and lower-cased; punctuation is retained.
// Frequencies are sorted by descending count, then ascending word.
func Count(r io.Reader) (Stats, error) {
	result := Stats{Frequencies: make([]Frequency, 0)}
	counts := make(map[string]int)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	for scanner.Scan() {
		result.Lines++
		for _, word := range strings.Fields(scanner.Text()) {
			counts[strings.ToLower(word)]++
			result.Words++
		}
	}
	if err := scanner.Err(); err != nil {
		return Stats{}, fmt.Errorf("scan text: %w", err)
	}
	result.Unique = len(counts)
	for word, count := range counts {
		result.Frequencies = append(result.Frequencies, Frequency{word, count})
	}
	slices.SortFunc(result.Frequencies, func(a, b Frequency) int {
		if order := cmp.Compare(b.Count, a.Count); order != 0 {
			return order
		}
		return cmp.Compare(a.Word, b.Word)
	})
	return result, nil
}
