package new

import (
	"mohazit/lang"
	"strconv"
	"strings"
)

type Statement struct {
	Keyword string
	KwToken *Token
	Args    []*Token
}

// toknes reads and returns the tokens of the next statement
func nextStmtTokens() []*Token {
	out := []*Token{}
	var t *Token
	for hasNextToken() {
		t = NextToken()
		if t == nil {
			continue
		}
		out = append(out, t)
		if t.Type == tLinefeed {
			return out
		}
	}
	return trimSpaceTokens(out)
}

// Next reads and returns the next statement in the input string
func NextStmt() (*Statement, error) {
	raw := nextStmtTokens()
	if len(raw) < 1 {
		return nil, nil
	}
	kwToken := raw[0]
	if kwToken.Type != tIdent {
		return nil, perrf(kwToken, "expected identifier, got %s", kwToken.Type.String())
	}
	kw := strings.ToLower(kwToken.Raw)
	args := []*Token{}
	for i := 1; i < len(raw); i++ {
		args = append(args, raw[i])
	}
	return &Statement{kw, kwToken, args}, nil
}

// Args reads a slice of objects from the given token slice
func parseObjectList(tkns []*Token) ([]*lang.Object, error) {
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
func parseObject(t []*Token) (*lang.Object, error) {
	t = trimSpaceTokens(t)
	switch t[0].Type {
	case tIdent, tUnknown:
		v := lang.NewStr(t[0].Raw)
		for i := 0; i < len(t); i++ {
			tkn := t[i]
			switch tkn.Type {
			case tIdent, tUnknown, tSpace:
				v.StrV += tkn.Raw
			default:
				return lang.NewNil(), perrf(tkn, "unexpected %s in string literal", tkn.Type.String())
			}
		}
		return v, nil
	case tLiteral:
		v, err := strconv.Atoi(t[0].Raw)
		return lang.NewInt(v), err
	default:
		return lang.NewNil(), perrf(t[0], "unexpected %s", t[0])
	}
}

// TrimSpace removes tSpace tokens from both ends of the given token slice
func trimSpaceTokens(t []*Token) []*Token {
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
