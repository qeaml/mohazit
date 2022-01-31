package lang

import (
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
func parseObjectList(tkns []*Token) ([]*Object, error) {
	out := []*Object{}
	ctx := ""
	for _, tkn := range tkns {
		switch tkn.Type {
		case tOper:
			if tkn.Raw == "\\" {
				out = append(out, NewStr(strings.TrimSpace(ctx)))
				ctx = ""
			} else {
				ctx += tkn.Raw
			}
		case tLiteral:
			if len(ctx) > 0 {
				out = append(out, NewStr(strings.TrimSpace(ctx)))
			}
			ctx = ""
			v, err := strconv.Atoi(tkn.Raw)
			if err != nil {
				return nil, err
			}
			out = append(out, NewInt(v))
		default:
			ctx += tkn.Raw
		}
	}
	if ctx != "" {
		out = append(out, NewStr(strings.TrimSpace(ctx)))
		ctx = ""
	}
	return out, nil
}

// Tokens2object reads a single object from the given token slice
func parseObject(t []*Token) (*Object, error) {
	t = trimSpaceTokens(t)
	switch t[0].Type {
	case tBracket:
		if t[0].Raw != "[" {
			return NewNil(), perrf(t[0], "expected [, got %s", t[0].Raw)
		}
		funcnames := []*Token{}
		argstart := 0
	funcLoop:
		for _, tkn := range t[1:] {
			argstart++
			switch tkn.Type {
			case tIdent:
				funcnames = append(funcnames, tkn)
			case tSpace:
				continue
			case tBracket:
				if tkn.Raw != "]" {
					return NewNil(), perrf(tkn, "expected ], got %s", tkn.Raw)
				}
				break funcLoop
			default:
				return NewNil(), perrf(tkn, "unexpected %s in function list", tkn.Type.String())
			}
		}
		funcfuncs := []VFunc{}
		for _, fn := range funcnames {
			if ff, ok := Funcs[strings.ToLower(fn.Raw)]; ok {
				funcfuncs = append(funcfuncs, ff)
			} else {
				return NewNil(), perrf(fn, "unknown function %s", fn.Raw)
			}
		}
		args, err := parseObjectList(trimSpaceTokens(t[argstart+1:]))
		if err != nil {
			return NewNil(), err
		}
		final, err := funcfuncs[0](args)
		if err != nil {
			return final, err
		}
		for _, f := range funcfuncs[1:] {
			final, err = f([]*Object{final})
			if err != nil {
				return final, err
			}
		}
		return final, nil
	case tIdent, tUnknown:
		v := NewStr(t[0].Raw)
		for i := 0; i < len(t); i++ {
			tkn := t[i]
			switch tkn.Type {
			case tIdent, tUnknown, tSpace:
				v.StrV += tkn.Raw
			default:
				return NewNil(), perrf(tkn, "unexpected %s in string literal", tkn.Type.String())
			}
		}
		return v, nil
	case tLiteral:
		v, err := strconv.Atoi(t[0].Raw)
		return NewInt(v), err
	default:
		return NewNil(), perrf(t[0], "unexpected %s", t[0])
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
