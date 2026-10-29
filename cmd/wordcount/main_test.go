package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/01reportplus/golang-tutorial/internal/wordcount"
)

func TestRunJSON(t *testing.T) {
	var out, diagnostics bytes.Buffer
	code := run([]string{"-json", "-text", "Go go 中文", "-top", "1"}, strings.NewReader("ignored"), &out, &diagnostics)
	if code != 0 {
		t.Fatalf("exit=%d: %s", code, &diagnostics)
	}
	var stats wordcount.Stats
	if err := json.Unmarshal(out.Bytes(), &stats); err != nil {
		t.Fatal(err)
	}
	if stats.Words != 3 || stats.Unique != 2 || len(stats.Frequencies) != 1 || stats.Frequencies[0].Word != "go" {
		t.Fatalf("unexpected result: %+v", stats)
	}
}

func TestRunInputAndErrors(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		code int
	}{
		{"stdin", nil, 0},
		{"empty explicit text", []string{"-text", ""}, 0},
		{"mutually exclusive", []string{"-text", "go", "-file", "x"}, 2},
		{"negative top", []string{"-top", "-1"}, 2},
		{"positional", []string{"unexpected"}, 2},
		{"missing file", []string{"-file", "does-not-exist.txt"}, 1},
		{"help", []string{"-h"}, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var out, diagnostics bytes.Buffer
			if got := run(tt.args, strings.NewReader("go"), &out, &diagnostics); got != tt.code {
				t.Fatalf("exit=%d, want %d; %s", got, tt.code, &diagnostics)
			}
		})
	}
}
