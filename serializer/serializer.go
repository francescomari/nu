package serializer

import (
	"fmt"
	"io"
	"strings"

	"github.com/francescomari/nu/parser"
)

// Serialize serializes a stream of commands into a io.Writer. If an error
// command is returned from the stream, or if an unexpected command is met,
// Serialize returns with a non-nil error.
func Serialize(commands <-chan parser.Cmd, w io.Writer) error {
	replacer := strings.NewReplacer("\n", "\\n", "\\", "\\\\")

	for command := range commands {
		switch cmd := command.(type) {
		case parser.R:
			fmt.Fprintf(w, "r\n")
		case parser.C:
			fmt.Fprintf(w, "c %v\n", cmd.Name)
		case parser.P:
			fmt.Fprintf(w, "p %v %v\n", cmd.Type, cmd.Name)
		case parser.V:
			fmt.Fprintf(w, "v %v\n", replacer.Replace(cmd.Data))
		case parser.X:
			fmt.Fprintf(w, "x %v\n", cmd.Data)
		case parser.Up:
			fmt.Fprintf(w, "^\n")
		case parser.Err:
			return fmt.Errorf("error at line %v: %v", cmd.Line, cmd.Err)
		default:
			return fmt.Errorf("unrecognized command: %#v", cmd)
		}
	}
	return nil
}
