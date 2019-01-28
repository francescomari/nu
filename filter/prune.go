package filter

import (
	"github.com/francescomari/nu/parser"
	"github.com/francescomari/nu/paths"
)

func Prune(path string, commands <-chan parser.Cmd) (<-chan parser.Cmd, error) {
	subtree, err := paths.Components(path)
	if err != nil {
		return nil, err
	}

	ch := make(chan parser.Cmd)
	go func() {
		defer close(ch)

		var (
			current []string
			emit    bool
		)

		for command := range commands {
			switch cmd := command.(type) {
			case parser.Err:
				ch <- cmd
			case parser.R:
				emit = !isInSubtree(current, subtree)
				if emit {
					ch <- cmd
				}
			case parser.C:
				current = append(current, cmd.Name)
				emit = !isInSubtree(current, subtree)
				if emit {
					ch <- cmd
				}
			case parser.P:
				current = append(current, cmd.Name)
				if emit {
					ch <- cmd
				}
			case parser.Up:
				if emit {
					ch <- cmd
				}
				if len(current) > 0 {
					current = current[:len(current)-1]
				}
				emit = !isInSubtree(current, subtree)
			default:
				if emit {
					ch <- cmd
				}
			}
		}
	}()
	return ch, nil
}
