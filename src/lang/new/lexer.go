package new

import (
	"fmt"
	"mohazit/lang"
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

// Lexer splits the input string into individual tokens
type Lexer struct {
	line   uint
	col    uint
	source string
	pos    int
}

// NewLexer creates an empty Lexer with an empty input string
func NewLexer() *Lexer {
	return &Lexer{1, 1, "", 0}
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
		return l.tkn(tSpace, toString(l.advance())), nil
	}

	if l.peek() == '\r' && l.peekNext() == '\n' {
		return l.tkn(tLinefeed, toString(l.advance())+toString(l.advance())), nil
	}
	if l.peek() == '\n' {
		return l.tkn(tLinefeed, toString(l.advance())), nil
	}

	if isBracket(l.peek()) {
		return l.tkn(tBracket, toString(l.advance())), nil
	}

	if isIdentStart(l.peek()) {
		ident := toString(l.advance())
		for isIdentCont(l.peek()) {
			ident += toString(l.advance())
		}
		return l.tkn(tIdent, ident), nil
	}

	if isDigit(l.peek()) || l.peek() == '-' {
		literal := toString(l.advance())
		for isDigit(l.peek()) {
			literal += toString(l.advance())
		}
		return l.tkn(tLiteral, literal), nil
	}

	if l.peek() == '\\' && l.peekNext() == ' ' {
		return l.tkn(tOper, toString(l.advance())), nil
	}

	dump := ""
	for !isSpace(l.peek()) && !isDigit(l.peek()) && l.peek() != '\r' && l.peek() != '\n' && l.peek() != 0 {
		dump += toString(l.advance())
	}
	for op := range lang.Comps {
		if dump == op {
			return l.tkn(tOper, dump), nil
		}
	}
	if len(dump) == 0 {
		return nil, nil
	}
	// return nil, badToken.Get(dump)
	return l.tkn(tUnknown, dump), nil
}

// Has returns true if there may be more tokens in the input string
func (l *Lexer) Has() bool {
	return l.pos != len(l.source)
}

func (l *Lexer) tkn(t TokenType, r string) *Token {
	token := &Token{l.line, l.col, t, r}
	if t == tLinefeed {
		l.line++
		l.col = 0
	}
	l.col += uint(len(r))
	return token
}
