package transform

import (
	"strings"

	"github.com/francescomari/nu/parser"
)

// PropertyPath is contains the fully qualified path and the type of a
// property, or an error if the transformation fails.
type PropertyPath struct {
	Path string
	Type string

	Err  error
	Line int
}

// Properties transforms a stream of commands into a stream of fully qualified
// paths of properties.
func Properties(cmds <-chan parser.Cmd) <-chan PropertyPath {
	results := make(chan PropertyPath)

	go func() {
		defer close(results)

		var path []string

		for cmd := range cmds {
			switch c := cmd.(type) {
			case parser.Err:
				results <- PropertyPath{Err: c.Err, Line: c.Line}
				break
			case parser.R:
				path = append(path, "")
			case parser.C:
				path = append(path, c.Name)
			case parser.P:
				path = append(path, c.Name)
				results <- PropertyPath{Path: strings.Join(path, "/"), Type: c.Type}
			case parser.Up:
				path = path[:len(path)-1]
			}
		}
	}()

	return results
}
