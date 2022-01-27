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

type Lexer struct {
	source string
	pos    int
}

func NewLexer() *Lexer {
	return &Lexer{"", 0}
}

func (l *Lexer) Source(src string) {
	l.source = src
	l.pos = 0
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.source) {
		return 0
	}
	return l.source[l.pos]
}

func (l *Lexer) peekNext() byte {
	if l.pos >= len(l.source)-1 {
		return 0
	}
	return l.source[l.pos+1]
}

func (l *Lexer) advance() byte {
	b := l.peek()
	l.pos++
	return b
}

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

func (l *Lexer) Has() bool {
	return l.pos != len(l.source)
}
