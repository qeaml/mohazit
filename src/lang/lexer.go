package lang

import (
	"fmt"
)

type TokenType uint8

const (
	tSpace TokenType = iota
	tLinefeed
	tLiteral
	tIdent
	tOper
	tBracket
	tRef
	tUnknown
)

func (t TokenType) String() string {
	switch t {
	case tSpace:
		return "space"
	case tLinefeed:
		return "linefeed"
	case tLiteral:
		return "literal"
	case tIdent:
		return "identifier"
	case tOper:
		return "operator"
	case tBracket:
		return "bracket"
	case tRef:
		return "ref"
	default:
		return "unknown"
	}
}

type Token struct {
	Line uint
	Col  uint
	Type TokenType
	Raw  string
}

func (t *Token) String() string {
	return fmt.Sprintf("<%s `%s` at %d:%d>", t.Type.String(), t.Raw, t.Line, t.Col)
}

var line uint = 1
var col uint = 1
var source string = ""
var pos int = 0

// Source sets the input string
func Source(src string) {
	source = src
	pos = 0
}

// peek returns the current character WITHOUT advancing the internal pointer.
// Returns 0 if there are no more readable characters.
func peek() byte {
	if pos >= len(source) {
		panic("excessive peek() call")
	}
	return source[pos]
}

// peekNext is the same as peek, but returns the next character over instead
func peekNext() byte {
	if pos+1 >= len(source) {
		panic("excessive peekNext() call")
	}
	return source[pos+1]
}

// advance is the same as peek, but DOES advance the internal pointer
func advance() byte {
	b := peek()
	pos++
	return b
}

// nextToken returns the next token in the input string, or nil if there are no
// more tokens left
func NextToken() *Token {
	c := peek()

	if isSpace(c) {
		return makeToken(tSpace, toString(advance()))
	}

	if c == '\r' && peekNext() == '\n' {
		return makeToken(tLinefeed, toString(advance())+toString(advance()))
	}
	if c == '\n' {
		return makeToken(tLinefeed, toString(advance()))
	}

	if isBracket(c) {
		if c == '{' {
			_ = advance()
			dump := ""
			for canAdvance() {
				if peek() == '}' {
					_ = advance()
					break
				}
				dump += toString(advance())
			}
			return makeToken(tRef, dump)
		}
		return makeToken(tBracket, toString(advance()))
	}

	if isIdentStart(c) {
		ident := toString(advance())
		for canAdvance() && isIdentCont(peek()) {
			ident += toString(advance())
		}
		return makeToken(tIdent, ident)
	}

	if isDigit(c) || c == '-' {
		literal := toString(advance())
		for canAdvance() && isDigit(peek()) {
			literal += toString(advance())
		}
		return makeToken(tLiteral, literal)
	}

	if c == '\\' {
		_ = advance()
		e := advance()
		switch e {
		case ' ':
			return makeTokenAlt(tOper, "\\", 2)
		case 'n':
			return makeTokenAlt(tUnknown, "\n", 2)
		case 'r':
			return makeTokenAlt(tUnknown, "\r", 2)
		case 't':
			return makeTokenAlt(tUnknown, "\t", 2)
		default:
			return makeToken(tUnknown, "\\"+toString(e))
		}
	}

	if isOper(c) {
		dump := toString(advance())
		for canAdvance() && isOper(peek()) {
			dump += toString(advance())
		}
		return makeToken(tOper, dump)
	}

	dump := toString(advance())
	for canAdvance() && !isValid(peek()) {
		dump += toString(advance())
	}
	return makeToken(tUnknown, dump)
}

// canAdvance returns true if there may be more tokens in the input string
func canAdvance() bool {
	return pos < len(source)
}

// makeToken creates a *Token from the input and advances the line and col
// counters
func makeToken(t TokenType, r string) *Token {
	token := &Token{line, col, t, r}
	if t == tLinefeed {
		line++
		col = 0
	}
	col += uint(len(r))
	return token
}

func makeTokenAlt(t TokenType, r string, len uint) *Token {
	token := &Token{line, col, t, r}
	if t == tLinefeed {
		line++
		col = 0
	}
	col += len
	return token
}

var OperChars []byte
