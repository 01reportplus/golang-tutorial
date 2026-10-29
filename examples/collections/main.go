package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	base := []int{1, 2, 3}
	copyOfBase := slices.Clone(base)
	copyOfBase[0] = 9
	fmt.Println("original:", base, "clone:", copyOfBase)
	counts := map[string]int{"go": 2, "中文": 1}
	for _, key := range slices.Sorted(maps.Keys(counts)) {
		fmt.Println(key, counts[key])
	}
}
