package new

import (
	"fmt"
	"mohazit/lang"
	"mohazit/lib"
	"strconv"
	"strings"
)

var (
	notEnough = lib.LazyError("parser: need more tokens (got %s, want 1)", "npar_notenough")
	notIdent  = lib.LazyError("parser: unexpected token (got %s, want 3)", "npar_notident")
	badKw     = lib.LazyError("parser: unknown keyword: %s", "npar_badkw")
)

type Statement struct {
	Keyword string
	Args    []*Token
}

type Parser struct {
	lexer *Lexer
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{lexer}
}

func (p *Parser) tokens() ([]*Token, error) {
	out := []*Token{}
	var t *Token
	var err error
	for p.lexer.Has() {
		t, err = p.lexer.Next()
		if err != nil {
			return nil, err
		}
		out = append(out, t)
		if t.Type == tLinefeed {
			return out, nil
		}
	}
	return out, nil
}

func (p *Parser) Next() (*Statement, error) {
	raw, err := p.tokens()
	if err != nil {
		return nil, err
	}
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
