package paths

import (
	"errors"
	"io"
	"strings"
)

var (
	// ErrInvalidPath is returned when a path is invalid.
	ErrInvalidPath = errors.New("invalid path")
)

// Components returns the components of a fully qualified path.
func Components(path string) ([]string, error) {
	const (
		stateStart = iota
		stateEnd
		stateSlash
		stateComponent
	)

	var (
		result    []string
		component string
		reader    = strings.NewReader(path)
		state     = stateStart
	)

	for {
		if state == stateEnd {
			return result, nil
		}

		c, _, err := reader.ReadRune()

		if err != nil && err != io.EOF {
			return nil, err
		}

		switch state {
		case stateStart:
			switch c {
			case 0:
				return nil, ErrInvalidPath
			case '/':
				state = stateSlash
			default:
				return nil, ErrInvalidPath
			}
		case stateSlash:
			switch c {
			case 0:
				state = stateEnd
			case '/':
				state = stateSlash
			default:
				component = string(c)
				state = stateComponent
			}
		case stateComponent:
			switch c {
			case 0:
				result = append(result, component)
				state = stateEnd
			case '/':
				result = append(result, component)
				state = stateSlash
			default:
				component += string(c)
				state = stateComponent
			}
		}
	}
}
