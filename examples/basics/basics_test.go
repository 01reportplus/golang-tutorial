package basics

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
	"testing/synctest"
	"time"
	"unicode/utf8"
)

func TestReverse(t *testing.T) {
	for _, tt := range []struct{ in, want string }{
		{"", ""}, {"Go中文", "文中oG"}, {"🙂a", "a🙂"},
	} {
		if got := Reverse(tt.in); got != tt.want {
			t.Fatalf("Reverse(%q)=%q, want %q", tt.in, got, tt.want)
		}
	}
}

func FuzzReverse(f *testing.F) {
	f.Add("Go中文")
	f.Add("")
	f.Fuzz(func(t *testing.T, text string) {
		if !utf8.ValidString(text) {
			t.Skip()
		}
		if got := Reverse(Reverse(text)); got != text {
			t.Fatalf("double reverse=%q, want %q", got, text)
		}
	})
}

func TestUniqueAndSum(t *testing.T) {
	type Score int
	if got := Sum([]Score{2, 3, 4}); got != 9 {
		t.Fatalf("sum=%d", got)
	}
	if got := Unique([]string{"a", "b", "a", "c", "b"}); !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Fatalf("unique=%v", got)
	}
}

func TestSliceAliasing(t *testing.T) {
	base := []int{1, 2, 3}
	shared := base[:1]
	shared = append(shared, 9)
	if base[1] != 9 {
		t.Fatal("expected shared backing storage")
	}
	isolated := base[:1:1]
	isolated = append(isolated, 7)
	if base[1] != 9 || isolated[1] != 7 {
		t.Fatal("full slice expression failed to isolate append")
	}
}

type problem struct{}

func (*problem) Error() string { return "problem" }

func TestTypedNil(t *testing.T) {
	var ptr *problem
	var err error = ptr
	if err == nil {
		t.Fatal("interface with dynamic type must not be nil")
	}
	var recovered *problem
	if !errors.As(err, &recovered) || recovered != nil {
		t.Fatal("typed nil dynamic value not preserved")
	}
}

func TestErrorWrapping(t *testing.T) {
	want := errors.New("not found")
	err := fmt.Errorf("load user: %w", want)
	if !errors.Is(err, want) {
		t.Fatal("lost error chain")
	}
}

func TestMapPresence(t *testing.T) {
	m := map[string]int{"present": 0}
	_, present := m["present"]
	_, absent := m["missing"]
	if !present || absent {
		t.Fatal("presence is different from the value zero")
	}
}

func deferredValues() (out []int) {
	n := 1
	defer func(value int) { out = append(out, value) }(n)
	defer func() { out = append(out, n) }()
	n = 2
	return
}

func TestDeferEvaluation(t *testing.T) {
	if got := deferredValues(); !slices.Equal(got, []int{2, 1}) {
		t.Fatalf("defer=%v", got)
	}
}

func TestLoopVariables(t *testing.T) {
	var funcs []func() int
	for i := range 3 {
		funcs = append(funcs, func() int { return i })
	}
	for i, fn := range funcs {
		if got := fn(); got != i {
			t.Fatalf("closure %d=%d", i, got)
		}
	}
}

func TestIteratorStops(t *testing.T) {
	var got []int
	for n := range Count(100) {
		got = append(got, n)
		if n == 2 {
			break
		}
	}
	if !slices.Equal(got, []int{0, 1, 2}) {
		t.Fatalf("iterator=%v", got)
	}
}

func TestContextDeadlineWithVirtualTime(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
		defer cancel()
		<-ctx.Done() // Virtual time advances without an actual one-hour wait.
		if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
			t.Fatalf("got %v", ctx.Err())
		}
	})
}

func TestNewExpression(t *testing.T) {
	p := new(42) // Go 1.26 expression form.
	if *p != 42 {
		t.Fatalf("new expression=%d", *p)
	}
}

func ExampleSum() {
	fmt.Println(Sum([]int{1, 2, 3}))
	// Output: 6
}
