package main

import (
	"fmt"

	"github.com/01reportplus/golang-tutorial/examples/basics"
)

func main() {
	fmt.Println(basics.Sum([]int64{10, 20, 30}))
	fmt.Println(basics.Unique([]string{"go", "java", "go"}))
	for value := range basics.Count(3) {
		fmt.Println(value)
	}
}
