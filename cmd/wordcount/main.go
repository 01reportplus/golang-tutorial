// Command wordcount counts whitespace-delimited tokens from text, a file, or stdin.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/01reportplus/golang-tutorial/internal/wordcount"
)

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("wordcount", flag.ContinueOnError)
	flags.SetOutput(stderr)
	text := flags.String("text", "", "text to count (exclusive with -file)")
	path := flags.String("file", "", "UTF-8 file to count")
	top := flags.Int("top", 10, "maximum frequencies; 0 means all")
	asJSON := flags.Bool("json", false, "output JSON")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	textProvided := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "text" {
			textProvided = true
		}
	})
	if flags.NArg() != 0 || *top < 0 || (textProvided && *path != "") {
		fmt.Fprintln(stderr, "use stdin, -text, or -file; -top must be nonnegative; positional arguments are not accepted")
		return 2
	}
	reader := stdin
	if textProvided {
		reader = strings.NewReader(*text)
	} else if *path != "" {
		file, err := os.Open(*path)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		defer file.Close()
		reader = file
	}
	stats, err := wordcount.Count(reader)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if *top > 0 && len(stats.Frequencies) > *top {
		stats.Frequencies = stats.Frequencies[:*top]
	}
	if *asJSON {
		err = json.NewEncoder(stdout).Encode(stats)
	} else {
		var output strings.Builder
		fmt.Fprintf(&output, "lines=%d words=%d unique=%d\n", stats.Lines, stats.Words, stats.Unique)
		for _, entry := range stats.Frequencies {
			fmt.Fprintf(&output, "%s\t%d\n", entry.Word, entry.Count)
		}
		_, err = io.WriteString(stdout, output.String())
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func main() { os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)) }
