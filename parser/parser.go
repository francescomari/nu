package parser

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// Parse parses an export from the specified io.Reader and emits a stream of
// commands.
func Parse(reader io.Reader) <-chan Cmd {
	ch := make(chan Cmd)
	go parse(reader, ch)
	return ch
}

func parse(reader io.Reader, ch chan<- Cmd) {
	defer close(ch)

	const (
		stateStart = iota
		stateEnd
		stateError
		stateR
		stateUp
		stateC
		stateCSpace
		stateCName
		stateP
		statePSpace
		statePSlash
		statePType
		statePTypeSpace
		statePName
		stateV
		stateVSpace
		stateVData
		stateVDataSlash
		stateX
		stateXSpace
		stateXData
	)

	var (
		cName    string
		pType    string
		pName    string
		data     strings.Builder
		state    = stateStart
		buffered = bufio.NewReader(reader)
		line     = 1
	)

	for {
		if state == stateEnd {
			break
		}

		if state == stateError {
			ch <- Err{Err: ErrInvalidInput, Line: line}
			break
		}

		c, _, err := buffered.ReadRune()

		if err != nil && err != io.EOF {
			ch <- Err{Err: err, Line: line}
			break
		}

		switch state {
		case stateStart:
			switch {
			case c == 0:
				state = stateEnd
			case c == '\n':
				line++
				state = stateStart
			case unicode.IsSpace(c):
				state = stateStart
			case c == 'r':
				state = stateR
			case c == '^':
				state = stateUp
			case c == 'c':
				state = stateC
			case c == 'p':
				state = stateP
			case c == 'v':
				state = stateV
			case c == 'x':
				state = stateX
			default:
				state = stateError
			}
		case stateR:
			switch {
			case c == 0:
				ch <- R{}
				state = stateEnd
			case c == '\n':
				line++
				ch <- R{}
				state = stateStart
			case unicode.IsSpace(c):
				state = stateR
			default:
				state = stateError
			}
		case stateUp:
			switch {
			case c == 0:
				ch <- Up{}
				state = stateEnd
			case c == '\n':
				line++
				ch <- Up{}
				state = stateStart
			case unicode.IsSpace(c):
				state = stateUp
			default:
				state = stateError
			}
		case stateC:
			switch {
			case unicode.IsSpace(c):
				state = stateCSpace
			default:
				state = stateError
			}
		case stateCSpace:
			switch {
			case c == 0:
				state = stateError
			case unicode.IsSpace(c):
				state = stateCSpace
			default:
				cName = string(c)
				state = stateCName
			}
		case stateCName:
			switch {
			case c == 0:
				ch <- C{cName}
				state = stateEnd
			case c == '\n':
				line++
				ch <- C{cName}
				state = stateStart
			default:
				cName += string(c)
				state = stateCName
			}
		case stateP:
			switch {
			case unicode.IsSpace(c):
				state = statePSpace
			default:
				state = stateError
			}
		case statePSpace:
			switch {
			case c == 0:
				state = stateError
			case c == '\n':
				state = stateError
			case unicode.IsSpace(c):
				state = statePSpace
			default:
				pType = string(c)
				state = statePType
			}
		case statePType:
			switch {
			case c == 0:
				state = stateError
			case c == '\n':
				state = stateError
			case unicode.IsSpace(c):
				state = statePTypeSpace
			default:
				pType += string(c)
				state = statePType
			}
		case statePTypeSpace:
			switch {
			case c == 0:
				state = stateError
			case c == '\n':
				state = stateError
			case unicode.IsSpace(c):
				state = statePTypeSpace
			default:
				pName = string(c)
				state = statePName
			}
		case statePName:
			switch {
			case c == 0:
				ch <- P{pType, pName}
				state = stateEnd
			case c == '\n':
				line++
				ch <- P{pType, pName}
				state = stateStart
			default:
				pName += string(c)
				state = statePName
			}
		case stateV:
			switch {
			case c == 0:
				ch <- V{}
				state = stateEnd
			case c == '\n':
				line++
				ch <- V{}
				state = stateStart
			case unicode.IsSpace(c):
				state = stateVSpace
			default:
				state = stateError
			}
		case stateVSpace:
			switch {
			case c == 0:
				ch <- V{}
				state = stateEnd
			case c == '\n':
				line++
				ch <- V{}
				state = stateStart
			case unicode.IsSpace(c):
				state = stateVSpace
			case c == '\\':
				data.Reset()
				state = stateVDataSlash
			default:
				data.Reset()
				data.WriteRune(c)
				state = stateVData
			}
		case stateVDataSlash:
			switch {
			case c == '\\':
				data.WriteRune('\\')
				state = stateVData
			case c == 'n':
				data.WriteRune('\n')
				state = stateVData
			default:
				state = stateError
			}
		case stateVData:
			switch {
			case c == 0:
				ch <- V{data.String()}
				state = stateEnd
			case c == '\n':
				line++
				ch <- V{data.String()}
				state = stateStart
			case c == '\\':
				state = stateVDataSlash
			default:
				data.WriteRune(c)
				state = stateVData
			}
		case stateX:
			switch {
			case c == 0:
				ch <- X{}
				state = stateEnd
			case c == '\n':
				line++
				ch <- X{}
				state = stateStart
			case unicode.IsSpace(c):
				state = stateXSpace
			default:
				state = stateError
			}
		case stateXSpace:
			switch {
			case c == 0:
				ch <- X{}
				state = stateEnd
			case c == '\n':
				line++
				ch <- X{}
				state = stateStart
			case unicode.IsSpace(c):
				state = stateXSpace
			default:
				data.Reset()
				data.WriteRune(c)
				state = stateXData
			}
		case stateXData:
			switch {
			case c == 0:
				ch <- X{data.String()}
				state = stateEnd
			case c == '\n':
				line++
				ch <- X{data.String()}
				state = stateStart
			default:
				data.WriteRune(c)
				state = stateXData
			}
		}
	}
}
