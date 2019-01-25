package transform

import (
	"strings"

	"github.com/francescomari/nu/parser"
)

// NodePath contains a fully qualified path of a node, or an error if the
// transformation fails.
type NodePath struct {
	Path string

	Err  error
	Line int
}

// Nodes transform a stream of commands into a stream of fully qualified node
// paths.
func Nodes(cmds <-chan parser.Cmd) <-chan NodePath {
	results := make(chan NodePath)

	go func() {
		defer close(results)

		var path []string

		for cmd := range cmds {
			switch c := cmd.(type) {
			case parser.Err:
				results <- NodePath{Err: c.Err, Line: c.Line}
				break
			case parser.R:
				path = append(path, "")
				results <- NodePath{Path: "/"}
			case parser.C:
				path = append(path, c.Name)
				results <- NodePath{Path: strings.Join(path, "/")}
			case parser.P:
				path = append(path, c.Name)
			case parser.Up:
				path = path[:len(path)-1]
			}
		}
	}()

	return results
}
