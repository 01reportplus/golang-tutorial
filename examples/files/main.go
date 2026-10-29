package main

import (
	"embed"
	"fmt"

	"github.com/01reportplus/golang-tutorial/internal/wordcount"
)

//go:embed message.txt
var files embed.FS

func main() {
	file, err := files.Open("message.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	stats, err := wordcount.Count(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("lines=%d words=%d\n", stats.Lines, stats.Words)
}
