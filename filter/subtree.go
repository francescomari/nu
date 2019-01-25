package filter

import (
	"fmt"

	"github.com/francescomari/nu/parser"
	"github.com/francescomari/nu/paths"
)

// Subtree filters a stream of commands into another stream of command where the
// tree rooted at `path` is the new root. Every part of the input commands not
// rooted at `path` is excluded from the output commands.
func Subtree(path string, commands <-chan parser.Cmd) (<-chan parser.Cmd, error) {
	subtree, err := paths.Components(path)
	if err != nil {
		return nil, fmt.Errorf("splitting path components: %v", err)
	}

	ch := make(chan parser.Cmd)
	go func() {
		defer close(ch)

		var (
			current []string
			send    bool
		)

		for command := range commands {
			switch cmd := command.(type) {
			case parser.Err:
				ch <- cmd
			case parser.R:
				send = isInSubtree(current, subtree)
				if send {
					ch <- cmd
				}
			case parser.C:
				current = append(current, cmd.Name)
				if send {
					ch <- cmd
					continue
				}
				send = isInSubtree(current, subtree)
				if send {
					ch <- parser.R{}
				}
			case parser.P:
				current = append(current, cmd.Name)
				if send {
					ch <- cmd
				}
			case parser.Up:
				if send {
					ch <- cmd
				}
				if len(current) > 0 {
					current = current[:len(current)-1]
				}
				send = isInSubtree(current, subtree)
			default:
				if send {
					ch <- cmd
				}
			}
		}
	}()
	return ch, nil
}

func isInSubtree(path, subtree []string) bool {
	if len(path) < len(subtree) {
		return false
	}
	for i := range subtree {
		if subtree[i] != path[i] {
			return false
		}
	}
	return true
}
