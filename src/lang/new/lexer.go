package new

import (
	"mohazit/lang"
)

type TokenType uint8

const (
	tSpace TokenType = iota
	tLinefeed
	tLiteral
	tIdent
	tOper
	tInvalid
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
		return "ident"
	case tOper:
		return "oper"
	default:
		return "invalid"
	}
}

type Token struct {
	Type TokenType
	Raw  string
}

// Lexer splits the input string into individual tokens
type Lexer struct {
	source string
	pos    int
}

// NewLexer creates an empty Lexer with an empty input string
func NewLexer() *Lexer {
	return &Lexer{"", 0}
}

// Source sets this Lexer's input string
func (l *Lexer) Source(src string) {
	l.source = src
	l.pos = 0
}

// peek returns the current character WITHOUT advancing the internal pointer.
// Returns 0 if there are no more readable characters.
func (l *Lexer) peek() byte {
	if l.pos >= len(l.source) {
		return 0
	}
	return l.source[l.pos]
}

// peekNext is the same as peek, but returns the next character over instead
func (l *Lexer) peekNext() byte {
	if l.pos >= len(l.source)-1 {
		return 0
	}
	return l.source[l.pos+1]
}

// advance is the same as peek, but DOES advance the internal pointer
func (l *Lexer) advance() byte {
	b := l.peek()
	l.pos++
	return b
}

// Next returns the next token in the input string, or nil if there are no
// more tokens left
func (l *Lexer) Next() (*Token, error) {
	if isSpace(l.peek()) {
		return &Token{tSpace, toString(l.advance())}, nil
	}

	if l.peek() == '\r' && l.peekNext() == '\n' {
		return &Token{tLinefeed, toString(l.advance()) + toString(l.advance())}, nil
	}
	if l.peek() == '\n' {
		return &Token{tLinefeed, toString(l.advance())}, nil
	}

	if isIdentStart(l.peek()) {
		ident := toString(l.advance())
		for isIdentCont(l.peek()) {
			ident += toString(l.advance())
		}
		return &Token{tIdent, ident}, nil
	}

	if isDigit(l.peek()) || l.peek() == '-' {
		literal := toString(l.advance())
		for isDigit(l.peek()) {
			literal += toString(l.advance())
		}
		return &Token{tLiteral, literal}, nil
	}

	if l.peek() == '\\' && l.peekNext() == ' ' {
		return &Token{tOper, toString(l.advance())}, nil
	}

	dump := ""
	for !isSpace(l.peek()) && l.peek() != '\r' && l.peek() != '\n' && l.peek() != 0 {
		dump += toString(l.advance())
	}
	for op := range lang.Comps {
		if dump == op {
			return &Token{tOper, dump}, nil
		}
	}
	if len(dump) == 0 {
		return nil, nil
	}
	// return nil, badToken.Get(dump)
	return &Token{tInvalid, dump}, nil
}

// Has returns true if there may be more tokens in the input string
func (l *Lexer) Has() bool {
	return l.pos != len(l.source)
}
