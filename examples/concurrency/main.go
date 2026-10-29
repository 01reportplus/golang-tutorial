package main

import (
	"fmt"
	"sync"
)

func main() {
	results := make([]int, 5)
	var wg sync.WaitGroup
	for i := range results {
		wg.Go(func() { results[i] = i * i })
	}
	wg.Wait()
	fmt.Println(results)
}
