package main

import (
	"bytes"
	"fmt"
	"io"
)

func greet(w io.Writer, name string) error {
	_, err := fmt.Fprintf(w, "hello, %s\n", name)
	return err
}

func main() {
	var buffer bytes.Buffer
	if err := greet(&buffer, "Go"); err != nil {
		panic(err)
	}
	fmt.Print(buffer.String())
}
