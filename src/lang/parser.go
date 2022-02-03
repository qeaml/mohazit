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
	raw := [][]*Token{}
outer:
	for i := 0; i < len(tkns); i++ {
		tkn := tkns[i]
		this := []*Token{}
		switch tkn.Type {
		case tSpace, tLinefeed:
			// ignore
		case tIdent, tUnknown, tBracket:
			for {
				this = append(this, tkn)
				if i+1 < len(tkns) {
					i++
					tkn = tkns[i]
					if tkn.Type == tOper && tkn.Raw == "\\" {
						raw = append(raw, this)
						continue outer
					}
				} else {
					raw = append(raw, this)
					break outer
				}
			}
		case tLiteral:
			raw = append(raw, []*Token{tkn})
		case tOper:
			return nil, perrf(tkn, "unexpected token: %s", tkn.Type)
		}
	}
	for _, src := range raw {
		o, err := parseObject(src)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, nil
}

// Tokens2object reads a single object from the given token slice
func parseObject(t []*Token) (*Object, error) {
	t = trimSpaceTokens(t)
	switch t[0].Type {
	case tRef:
		if len(t) > 1 {
			return nil, perrf(t[1], "unexpected %s in reference", t[1].Type)
		}
		lv, ok := GetLocalVar(t[0].Raw)
		if ok {
			return lv, nil
		}
		gv, ok := GetGlobalVar(t[0].Raw)
		if ok {
			return gv, nil
		}
		return nil, perrf(t[0], "could not find variable %s", t[0].Raw)
	case tBracket:
		if t[0].Raw != "[" {
			// return NewNil(), perrf(t[0], "expected [, got %s", t[0].Raw)
			v := NewStr(t[0].Raw)
			for _, tkn := range t[1:] {
				switch tkn.Type {
				case tIdent, tUnknown, tSpace, tLiteral, tBracket:
					v.StrV += tkn.Raw
				default:
					return NewNil(), perrf(tkn, "unexpected %s in string literal", tkn.Type.String())
				}
			}
			return v, nil
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
		for _, tkn := range t[1:] {
			switch tkn.Type {
			case tIdent, tUnknown, tSpace, tLiteral, tBracket:
				v.StrV += tkn.Raw
			default:
				return NewNil(), perrf(tkn, "unexpected %s in string literal", tkn.Type.String())
			}
		}
		br := strings.TrimSpace(strings.ToLower(v.StrV))
		if br == "true" || br == "yes" {
			return NewBool(true), nil
		} else if br == "false" || br == "no" {
			return NewBool(false), nil
		} else if br == "nil" {
			return NewNil(), nil
		}
		return v, nil
	case tLiteral:
		v, err := strconv.Atoi(t[0].Raw)
		return NewInt(v), err
	default:
		return NewNil(), perrf(t[0], "unexpected %s in object value", t[0])
	}
}

// TrimSpace removes tSpace tokens from both ends of the given token slice
func trimSpaceTokens(t []*Token) []*Token {
	ltrim := []*Token{}
	ignore := true
	for _, tkn := range t {
		if (tkn.Type != tSpace && tkn.Type != tLinefeed) && ignore {
			ignore = false
		}
		if !ignore {
			ltrim = append(ltrim, tkn)
		}
	}
	rtrim := []*Token{}
	ignore = true
	for i := len(ltrim) - 1; i >= 0; i-- {
		if (ltrim[i].Type != tSpace && ltrim[i].Type != tLinefeed) && ignore {
			ignore = false
		}
		if !ignore {
			rtrim = append([]*Token{ltrim[i]}, rtrim...)
		}
	}
	return rtrim
}
