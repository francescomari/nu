package filter

import (
	"testing"

	"github.com/francescomari/nu/parser"
)

func TestPruneRoot(t *testing.T) {
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

	outCh, err := Prune("/", inCh)
	if err != nil {
		t.Fatalf("Prune: %v\n", err)
	}

	var out []parser.Cmd
	for cmd := range outCh {
		out = append(out, cmd)
	}

	assertCommandsEqual(t, nil, out)
}

func TestPruneNonExistentPath(t *testing.T) {
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

	outCh, err := Prune("/this/does/not/exist", inCh)
	if err != nil {
		t.Fatalf("Prune: %v\n", err)
	}

	var out []parser.Cmd
	for cmd := range outCh {
		out = append(out, cmd)
	}

	assertCommandsEqual(t, in, out)
}

func TestPruneLevelOne(t *testing.T) {
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

	outCh, err := Prune("/a", inCh)
	if err != nil {
		t.Fatalf("Prune: %v\n", err)
	}

	var out []parser.Cmd
	for cmd := range outCh {
		out = append(out, cmd)
	}

	expect := []parser.Cmd{
		parser.R{},
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

	assertCommandsEqual(t, expect, out)
}

func TestPruneLevelTwo(t *testing.T) {
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

	outCh, err := Prune("/a/c", inCh)
	if err != nil {
		t.Fatalf("Prune: %v\n", err)
	}

	var out []parser.Cmd
	for cmd := range outCh {
		out = append(out, cmd)
	}

	expect := []parser.Cmd{
		parser.R{},
		parser.C{Name: "a"},
		parser.P{Type: "string", Name: "p"},
		parser.V{Data: "a"},
		parser.Up{}, // End of /a[p]
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

	assertCommandsEqual(t, expect, out)
}
