package parser

// Cmd is the interface implemented by the commands returned by the parser.
type Cmd interface {
	cmd()
}

// R is the `r` command from an export.
type R struct {
}

func (R) cmd() {
}

// Up is the `^` command from an export.
type Up struct {
}

func (Up) cmd() {
}

// C is the `c` command from an export.
type C struct {
	// Name is the name of the node.
	Name string
}

func (C) cmd() {
}

// P is the `p` command from an export.
type P struct {
	// Type is the type of the property.
	Type string
	// Name is the name of the property.
	Name string
}

func (P) cmd() {
}

// V is the `v` command from an export.
type V struct {
	Data string
}

func (V) cmd() {
}

// X is the `x` command from an export.
type X struct {
	Data string
}

func (X) cmd() {
}

// Err is an error emitted from the parser.
type Err struct {
	// Err is the error itself.
	Err error
	// Line is the line where the error was detected.
	Line int
}

func (Err) cmd() {
}
