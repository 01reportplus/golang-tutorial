// Package basics contains small algorithms used to explore Go semantics.
package basics

import "iter"

type Integer interface{ ~int | ~int32 | ~int64 }

func Sum[T Integer](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

// Unique returns the first occurrence of each value, preserving order.
func Unique[T comparable](xs []T) []T {
	seen := make(map[T]struct{}, len(xs))
	out := make([]T, 0, len(xs))
	for _, x := range xs {
		if _, ok := seen[x]; ok {
			continue
		}
		seen[x] = struct{}{}
		out = append(out, x)
	}
	return out
}

// Reverse reverses Unicode code points, not grapheme clusters.
// Invalid UTF-8 is replaced during conversion to []rune.
func Reverse(s string) string {
	runes := []rune(s)
	for left, right := 0, len(runes)-1; left < right; left, right = left+1, right-1 {
		runes[left], runes[right] = runes[right], runes[left]
	}
	return string(runes)
}

// Count yields [0, n), stopping immediately when the consumer stops.
func Count(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < n; i++ {
			if !yield(i) {
				return
			}
		}
	}
}
