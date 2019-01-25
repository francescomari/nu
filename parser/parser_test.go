package parser

import (
	"strings"
	"testing"
)

func parseAll(s string) []Cmd {
	var cmds []Cmd
	for c := range Parse(strings.NewReader(s)) {
		cmds = append(cmds, c)
	}
	return cmds
}

func TestParseCmd(t *testing.T) {
	tests := []struct {
		line     string
		expected Cmd
	}{
		{"r", R{}},
		{"^", Up{}},
		{"c name", C{"name"}},
		{"p type name", P{"type", "name"}},
		{"v data with spaces", V{"data with spaces"}},
		{"v data\\nwith\\nnewlines", V{"data\nwith\nnewlines"}},
		{"v data\\\\with\\\\slashes", V{"data\\with\\slashes"}},
		{" r ", R{}},
		{" ^ ", Up{}},
		{" c  name", C{"name"}},
		{" p  type  name", P{"type", "name"}},
		{" v  data", V{"data"}},
		{"v \\ndata", V{"\ndata"}},
		{"v \\\\data", V{"\\data"}},
		{"x", X{}},
		{"x 0123456789abcdef", X{"0123456789abcdef"}},
		{"x 0123456789ABCDEF", X{"0123456789ABCDEF"}},
		{" x  F0", X{"F0"}},
		{"v ", V{}},
	}

	for _, tt := range tests {
		all := parseAll(tt.line)
		if len(all) != 1 {
			t.Errorf("parsing '%v': expected 1 command, got %v\n", tt.line, len(all))
			continue
		}
		cmd := all[0]
		if cmd != tt.expected {
			t.Errorf("parsing '%v': expected %v, got %v\n", tt.line, tt.expected, cmd)
		}
	}
}

func TestParse(t *testing.T) {
	cmds := `
	r
	c foo
	p string bar
	v baz
	^
	p binary foo
	x deadbeef
	^
	^
	^
	`
	expected := []Cmd{
		R{},
		C{"foo"},
		P{"string", "bar"},
		V{"baz"},
		Up{},
		P{"binary", "foo"},
		X{"deadbeef"},
		Up{},
		Up{},
		Up{},
	}
	all := parseAll(cmds)
	if len(all) != len(expected) {
		t.Fatalf("expcted %d commands, got %d\n", len(expected), len(all))
	}
	for i, a := range expected {
		if a != all[i] {
			t.Errorf("expected %v, got %v\n", a, all[i])
		}
	}
}

func TestParseVData(t *testing.T) {
	tests := []struct {
		line string
		data string
	}{
		{"v", ""},
		{"v\n", ""},
		{"v ", ""},
		{"v \n", ""},
		{" v", ""},
		{"v data", "data"},
		{"v data with spaces", "data with spaces"},
		{"v data\\\\with\\\\slashes", "data\\with\\slashes"},
		{"v data\\nwith\\nnewlines", "data\nwith\nnewlines"},
	}

	for _, tt := range tests {
		cmds := parseAll(tt.line)
		if len(cmds) != 1 {
			t.Errorf("parsing '%v': expected 1 command, got %v\n", tt.line, len(cmds))
			continue
		}

		v, ok := cmds[0].(V)
		if !ok {
			t.Errorf("parsing '%v': expected V, got %T\n", tt.line, cmds[0])
			continue
		}
		if v.Data != tt.data {
			t.Errorf("parsing '%v': expected data '%v', got '%v'\n", tt.line, tt.data, v.Data)
			continue
		}
	}
}

func TestParseXData(t *testing.T) {
	tests := []struct {
		line string
		data string
	}{
		{"x", ""},
		{"x\n", ""},
		{"x ", ""},
		{"x \n", ""},
		{" x", ""},
		{"x data", "data"},
	}

	for _, tt := range tests {
		cmds := parseAll(tt.line)
		if len(cmds) != 1 {
			t.Errorf("parsing '%v': expected 1 command, got %v\n", tt.line, len(cmds))
			continue
		}

		x, ok := cmds[0].(X)
		if !ok {
			t.Errorf("parsing '%v': expected X, got %T\n", tt.line, cmds[0])
			continue
		}
		if x.Data != tt.data {
			t.Errorf("parsing '%v': expected data '%v', got '%v'\n", tt.line, tt.data, x.Data)
			continue
		}
	}
}
