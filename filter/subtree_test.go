package filter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/francescomari/nu/parser"
)

func TestUnchangedSubtree(t *testing.T) {
	in := []parser.Cmd{
		parser.R{},
		parser.C{Name: "a"},
		parser.P{Type: "string", Name: "p"},
		parser.V{Data: "a"},
		parser.Up{}, // End of /a[p]
		parser.C{Name: "c"},
		parser.P{Type: "string", Name: "q"},
		parser.V{Data: "c"},
		parser.Up{}, // End of /a/c[q]
		parser.Up{}, // End of /a/c
		parser.Up{}, // End of /a
		parser.C{Name: "b"},
		parser.P{Type: "string", Name: "r"},
		parser.V{Data: "b"},
		parser.Up{}, // End of /b[r]
		parser.C{Name: "d"},
		parser.P{Type: "string", Name: "s"},
		parser.V{Data: "d"},
		parser.Up{}, // End of /b/d[s]
		parser.Up{}, // End of /b/d
		parser.Up{}, // End of /b
		parser.Up{}, // End of /
	}

	inCh := make(chan parser.Cmd)

	go func() {
		defer close(inCh)
		for _, cmd := range in {
			inCh <- cmd
		}
	}()

	outCh, err := Subtree("/", inCh)
	if err != nil {
		t.Fatalf("Subtree: %v\n", err)
	}

	var out []parser.Cmd
	for cmd := range outCh {
		out = append(out, cmd)
	}

	assertCommandsEqual(t, in, out)
}

func TestLevelOneSubtree(t *testing.T) {
	in := []parser.Cmd{
		parser.R{},
		parser.C{Name: "a"},
		parser.P{Type: "string", Name: "p"},
		parser.V{Data: "a"},
		parser.Up{}, // End of /a[p]
		parser.C{Name: "c"},
		parser.P{Type: "string", Name: "q"},
		parser.V{Data: "c"},
		parser.Up{}, // End of /a/c[q]
		parser.Up{}, // End of /a/c
		parser.Up{}, // End of /a
		parser.C{Name: "b"},
		parser.P{Type: "string", Name: "r"},
		parser.V{Data: "b"},
		parser.Up{}, // End of /b[r]
		parser.C{Name: "d"},
		parser.P{Type: "string", Name: "s"},
		parser.V{Data: "d"},
		parser.Up{}, // End of /b/d[s]
		parser.Up{}, // End of /b/d
		parser.Up{}, // End of /b
		parser.Up{}, // End of /
	}

	inCh := make(chan parser.Cmd)

	go func() {
		defer close(inCh)
		for _, cmd := range in {
			inCh <- cmd
		}
	}()

	outCh, err := Subtree("/a", inCh)
	if err != nil {
		t.Fatalf("Subtree: %v\n", err)
	}

	var out []parser.Cmd
	for cmd := range outCh {
		out = append(out, cmd)
	}

	expect := []parser.Cmd{
		parser.R{},
		parser.P{Type: "string", Name: "p"},
		parser.V{Data: "a"},
		parser.Up{}, // End of /[p]
		parser.C{Name: "c"},
		parser.P{Type: "string", Name: "q"},
		parser.V{Data: "c"},
		parser.Up{}, // End of /c[q]
		parser.Up{}, // End of /c
		parser.Up{}, // End of /
	}

	assertCommandsEqual(t, expect, out)
}

func TestLevelTwoSubtree(t *testing.T) {
	in := []parser.Cmd{
		parser.R{},
		parser.C{Name: "a"},
		parser.P{Type: "string", Name: "p"},
		parser.V{Data: "a"},
		parser.Up{}, // End of /a[p]
		parser.C{Name: "c"},
		parser.P{Type: "string", Name: "q"},
		parser.V{Data: "c"},
		parser.Up{}, // End of /a/c[q]
		parser.Up{}, // End of /a/c
		parser.Up{}, // End of /a
		parser.C{Name: "b"},
		parser.P{Type: "string", Name: "r"},
		parser.V{Data: "b"},
		parser.Up{}, // End of /b[r]
		parser.C{Name: "d"},
		parser.P{Type: "string", Name: "s"},
		parser.V{Data: "d"},
		parser.Up{}, // End of /b/d[s]
		parser.Up{}, // End of /b/d
		parser.Up{}, // End of /b
		parser.Up{}, // End of /
	}

	inCh := make(chan parser.Cmd)

	go func() {
		defer close(inCh)
		for _, cmd := range in {
			inCh <- cmd
		}
	}()

	outCh, err := Subtree("/a/c", inCh)
	if err != nil {
		t.Fatalf("Subtree: %v\n", err)
	}

	var out []parser.Cmd
	for cmd := range outCh {
		out = append(out, cmd)
	}

	expect := []parser.Cmd{
		parser.R{},
		parser.P{Type: "string", Name: "q"},
		parser.V{Data: "c"},
		parser.Up{}, // End of /[q]
		parser.Up{}, // End of /
	}

	assertCommandsEqual(t, expect, out)
}

func TestNonExistentSubtree(t *testing.T) {
	in := []parser.Cmd{
		parser.R{},
		parser.C{Name: "a"},
		parser.P{Type: "string", Name: "p"},
		parser.V{Data: "a"},
		parser.Up{}, // End of /a[p]
		parser.C{Name: "c"},
		parser.P{Type: "string", Name: "q"},
		parser.V{Data: "c"},
		parser.Up{}, // End of /a/c[q]
		parser.Up{}, // End of /a/c
		parser.Up{}, // End of /a
		parser.C{Name: "b"},
		parser.P{Type: "string", Name: "r"},
		parser.V{Data: "b"},
		parser.Up{}, // End of /b[r]
		parser.C{Name: "d"},
		parser.P{Type: "string", Name: "s"},
		parser.V{Data: "d"},
		parser.Up{}, // End of /b/d[s]
		parser.Up{}, // End of /b/d
		parser.Up{}, // End of /b
		parser.Up{}, // End of /
	}

	inCh := make(chan parser.Cmd)

	go func() {
		defer close(inCh)
		for _, cmd := range in {
			inCh <- cmd
		}
	}()

	outCh, err := Subtree("/nope", inCh)
	if err != nil {
		t.Fatalf("Subtree: %v\n", err)
	}

	var out []parser.Cmd
	for cmd := range outCh {
		out = append(out, cmd)
	}

	assertCommandsEqual(t, nil, out)
}

func commandsEqual(a, b []parser.Cmd) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func assertCommandsEqual(t *testing.T, expected, got []parser.Cmd) {
	t.Helper()
	if commandsEqual(expected, got) {
		return
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Commands don't match.\n")
	fmt.Fprintf(&b, "Expected commands:\n")
	for _, cmd := range expected {
		fmt.Fprintf(&b, "  %#v\n", cmd)
	}
	fmt.Fprintf(&b, "Received commands:\n")
	for _, cmd := range got {
		fmt.Fprintf(&b, "  %#v\n", cmd)
	}
	t.Error(b.String())
}
