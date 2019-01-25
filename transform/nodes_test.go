package transform

import (
	"strings"
	"testing"

	"github.com/francescomari/nu/parser"
)

func TestNodes(t *testing.T) {
	cmds := parser.Parse(strings.NewReader(`
		r
		c 1
		p t n
		v x
		^
		c 1.1
		p t n
		v x
		^
		^
		c 1.2
		^
		^
		^
	`))

	var paths []string

	for p := range Nodes(cmds) {
		if p.Err != nil {
			t.Fatalf("error at line %v: %v\n", p.Line, p.Err)
		}
		paths = append(paths, p.Path)
	}

	expected := []string{
		"/",
		"/1",
		"/1/1.1",
		"/1/1.2",
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
