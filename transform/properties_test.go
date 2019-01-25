package transform

import (
	"strings"
	"testing"

	"github.com/francescomari/nu/parser"
)

func TestProperties(t *testing.T) {
	cmds := parser.Parse(strings.NewReader(`
		r
		p x y
		^
		c 1
		p t u
		^
		c 1.1
		p v w
		^
		^
		^
		^
	`))

	var paths []PropertyPath

	for p := range Properties(cmds) {
		if p.Err != nil {
			t.Fatalf("error at line %v: %v\n", p.Line, p.Err)
		}
		paths = append(paths, p)
	}

	expected := []PropertyPath{
		{Path: "/y", Type: "x"},
		{Path: "/1/u", Type: "t"},
		{Path: "/1/1.1/w", Type: "v"},
	}

	if len(paths) != len(expected) {
		t.Fatalf("expected %v paths, got %v\n", len(expected), len(paths))
	}
	for i, p := range expected {
		if p != paths[i] {
			t.Errorf("expected %v, got %v\n", p, paths[i])
		}
	}
}
