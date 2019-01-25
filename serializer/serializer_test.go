package serializer

import (
	"strings"
	"testing"

	"github.com/francescomari/nu/parser"
)

func TestSerializer(t *testing.T) {
	ch := make(chan parser.Cmd)

	go func() {
		defer close(ch)
		for _, cmd := range []parser.Cmd{
			parser.R{},
			parser.C{Name: "a"},
			parser.P{Type: "b", Name: "c"},
			parser.X{Data: "d"},
			parser.V{Data: "e"},
			parser.V{Data: "a\\b\\"},
			parser.V{Data: "a\nb\n"},
			parser.Up{},
		} {
			ch <- cmd
		}
	}()

	var w strings.Builder

	if err := Serialize(ch, &w); err != nil {
		t.Fatalf("serialize: %v\n", err)
	}

	var e string

	e += "r\n"
	e += "c a\n"
	e += "p b c\n"
	e += "x d\n"
	e += "v e\n"
	e += "v a\\\\b\\\\\n"
	e += "v a\\nb\\n\n"
	e += "^\n"

	if e != w.String() {
		t.Fatalf("unexpected output:\n%v", w.String())
	}
}
