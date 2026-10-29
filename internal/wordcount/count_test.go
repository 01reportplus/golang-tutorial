package wordcount

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestCount(t *testing.T) {
	for _, tt := range []struct {
		name, text           string
		lines, words, unique int
	}{
		{"empty", "", 0, 0, 0},
		{"blank lines", "\n\n", 2, 0, 0},
		{"trailing newline", "go\n", 1, 1, 1},
		{"case and Unicode whitespace", "Go go\n中文\u3000Go", 2, 4, 2},
		{"punctuation retained", "go Go!", 1, 2, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Count(strings.NewReader(tt.text))
			if err != nil {
				t.Fatal(err)
			}
			if got.Lines != tt.lines || got.Words != tt.words || got.Unique != tt.unique {
				t.Fatalf("got %+v, want lines=%d words=%d unique=%d", got, tt.lines, tt.words, tt.unique)
			}
		})
	}
}

func TestFrequencyOrder(t *testing.T) {
	got, err := Count(strings.NewReader("b a b a c"))
	if err != nil {
		t.Fatal(err)
	}
	want := []Frequency{{"a", 2}, {"b", 2}, {"c", 1}}
	for i, entry := range want {
		if got.Frequencies[i] != entry {
			t.Fatalf("entry %d: got %+v, want %+v", i, got.Frequencies[i], entry)
		}
	}
}

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

func TestReadError(t *testing.T) {
	want := errors.New("broken input")
	_, err := Count(io.MultiReader(strings.NewReader("go\n"), failingReader{want}))
	if !errors.Is(err, want) {
		t.Fatalf("got %v, want wrapped %v", err, want)
	}
}

func TestLongLine(t *testing.T) {
	if _, err := Count(strings.NewReader(strings.Repeat("x", 1024*1024+1))); err == nil {
		t.Fatal("expected scanner limit error")
	}
}

func BenchmarkCount(b *testing.B) {
	text := strings.Repeat("Go go tutorial 中文\n", 100)
	for b.Loop() {
		if _, err := Count(strings.NewReader(text)); err != nil {
			b.Fatal(err)
		}
	}
}
