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
		return 0
	}
	return source[pos]
}

// peekNext is the same as peek, but returns the next character over instead
func peekNext() byte {
	if pos >= len(source)-1 {
		return 0
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
	if isSpace(peek()) {
		return makeToken(tSpace, toString(advance()))
	}

	if peek() == '\r' && peekNext() == '\n' {
		return makeToken(tLinefeed, toString(advance())+toString(advance()))
	}
	if peek() == '\n' {
		return makeToken(tLinefeed, toString(advance()))
	}

	if isBracket(peek()) {
		return makeToken(tBracket, toString(advance()))
	}

	if isIdentStart(peek()) {
		ident := toString(advance())
		for isIdentCont(peek()) {
			ident += toString(advance())
		}
		return makeToken(tIdent, ident)
	}

	if isDigit(peek()) || peek() == '-' {
		literal := toString(advance())
		for isDigit(peek()) {
			literal += toString(advance())
		}
		return makeToken(tLiteral, literal)
	}

	if peek() == '\\' {
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

	dump := ""
	for !isSpace(peek()) && !isDigit(peek()) && !isBracket(peek()) && peek() != '\r' && peek() != '\n' && peek() != 0 {
		dump += toString(advance())
	}
	for op := range Comps {
		if dump == op {
			return makeToken(tOper, dump)
		}
	}
	if len(dump) == 0 {
		return nil
	}
	return makeToken(tUnknown, dump)
}

// hasNextToken returns true if there may be more tokens in the input string
func hasNextToken() bool {
	return pos != len(source)
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
