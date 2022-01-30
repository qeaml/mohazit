package new

import (
	"fmt"
	"mohazit/lang"
	"mohazit/lib"
	"strconv"
	"strings"
)

var (
	notIdent = lib.LazyError("parser: unexpected token (got %s, want ident)", "npar_notident")
)

type Statement struct {
	Keyword string
	Args    []*Token
}

// Parser reads it's Lexer's tokens and turns them into individual statements
type Parser struct {
	lexer *Lexer
}

// NewParser creates an empty Parser with an empty input string
func NewParser() *Parser {
	return &Parser{NewLexer()}
}

// Source sets this Parser's input string
func (p *Parser) Source(src string) {
	p.lexer.Source(src)
}

// toknes reads and returns the tokens of the next statement
func (p *Parser) tokens() []*Token {
	out := []*Token{}
	var t *Token
	for p.lexer.Has() {
		t = p.lexer.Next()
		if t == nil {
			continue
		}
		out = append(out, t)
		if t.Type == tLinefeed {
			return out
		}
	}
	return p.TrimSpace(out)
}

// Next reads and returns the next statement in the input string
func (p *Parser) Next() (*Statement, error) {
	raw := p.tokens()
	if len(raw) < 1 {
		return nil, nil
	}
	kwToken := raw[0]
	if kwToken.Type != tIdent {
		return nil, notIdent.Get(fmt.Sprint(kwToken.Type))
	}
	kw := strings.ToLower(kwToken.Raw)
	args := []*Token{}
	for i := 1; i < len(raw); i++ {
		args = append(args, raw[i])
	}
	return &Statement{kw, args}, nil
}

// Args reads a slice of objects from the given token slice
func (p *Parser) Args(tkns []*Token) ([]*lang.Object, error) {
	out := []*lang.Object{}
	ctx := ""
	for _, tkn := range tkns {
		switch tkn.Type {
		case tOper:
			if tkn.Raw == "\\" {
				out = append(out, lang.NewStr(strings.TrimSpace(ctx)))
				ctx = ""
			} else {
				ctx += tkn.Raw
			}
		case tLiteral:
			out = append(out, lang.NewStr(strings.TrimSpace(ctx)))
			ctx = ""
			v, err := strconv.Atoi(tkn.Raw)
			if err != nil {
				return nil, err
			}
			out = append(out, lang.NewInt(v))
		default:
			ctx += tkn.Raw
		}
	}
	if ctx != "" {
		out = append(out, lang.NewStr(strings.TrimSpace(ctx)))
		ctx = ""
	}
	return out, nil
}

// Tokens2object reads a single object from the given token slice
func (p *Parser) Tokens2object(t []*Token) (*lang.Object, error) {
	t = p.TrimSpace(t)
	switch t[0].Type {
	case tIdent, tUnknown:
		v := lang.NewStr(t[0].Raw)
		for i := 0; i < len(t); i++ {
			tkn := t[i]
			switch tkn.Type {
			case tIdent, tUnknown, tSpace:
				v.StrV += tkn.Raw
			default:
				return lang.NewNil(), unexTkn.Get(tkn.Type.String())
			}
		}
		return v, nil
	case tLiteral:
		v, err := strconv.Atoi(t[0].Raw)
		return lang.NewInt(v), err
	default:
		return lang.NewNil(), unexTkn.Get(t[0].Type.String())
	}
}

// TrimSpace removes tSpace tokens from both ends of the given token slice
func (p *Parser) TrimSpace(t []*Token) []*Token {
	ltrim := []*Token{}
	ignore := true
	for _, tkn := range t {
		if tkn.Type != tSpace && ignore {
			ignore = false
		}
		if !ignore {
			ltrim = append(ltrim, tkn)
		}
	}
	rtrim := []*Token{}
	ignore = true
	for i := len(ltrim) - 1; i >= 0; i-- {
		if ltrim[i].Type != tSpace && ignore {
			ignore = false
		}
		if !ignore {
			rtrim = append([]*Token{ltrim[i]}, rtrim...)
		}
	}
	return rtrim
}
