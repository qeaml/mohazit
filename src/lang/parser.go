package lang

import (
	"errors"
	"strconv"
	"strings"
)

type StatementType uint8

const (
	stCall StatementType = iota
	stSet
	stIf
	stElse
	stLabel
)

type Statement struct {
	Type StatementType
	Func string
	Args []*Object
}

func (s *Statement) Repr() string {
	strArgs := []string{}
	for _, a := range s.Args {
		strArgs = append(strArgs, a.Repr())
	}
	return s.Func + "(" + strings.Join(strArgs, "; ") + ")"
}

type Parser struct {
}

func (p *Parser) isWhitespace(c rune) bool {
	return c == ' ' || c == '\t'
}

func (p *Parser) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (p *Parser) id(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func (p *Parser) objNil() *Object {
	return &Object{
		Type: objNil,
	}
}

func (p *Parser) objInt(i int) *Object {
	return &Object{
		Type: objInt,
		IntV: i,
	}
}

func (p *Parser) objBool(b bool) *Object {
	return &Object{
		Type:  objBool,
		BoolV: b,
	}
}

func (p *Parser) objStr(s string) *Object {
	return &Object{
		Type: objStr,
		StrV: s,
	}
}

func (p *Parser) typeOf(s string) ObjectType {
	s = strings.TrimSpace(s)
	switch strings.ToLower(s) {
	case "nil":
		return objNil
	case "true", "yes", "false", "no":
		return objBool
	}
	if s[0] == '-' || p.isDigit(s[0]) {
		return objInt
	}
	return objStr
}

func (p *Parser) parseObject(s string) (*Object, error) {
	s = strings.TrimSpace(s)
	t := p.typeOf(s)
	switch t {
	case objNil:
		return p.objNil(), nil
	case objBool:
		s = strings.ToLower(s)
		if s == "true" || s == "yes" {
			return p.objBool(true), nil
		}
		if s == "false" || s == "no" {
			return p.objBool(false), nil
		}
		return p.objNil(), errors.New("invalid boolean value: " + s)
	case objInt:
		i, err := strconv.Atoi(s)
		if err != nil {
			return p.objNil(), err
		}
		return p.objInt(i), nil
	case objStr:
		return p.objStr(s), nil
	}
	return p.objNil(), errors.New("could not deterime type of value: " + s)
}

func (p *Parser) parseArgs(a []string) ([]*Object, error) {
	if len(a) == 0 {
		return []*Object{}, nil
	}
	if p.typeOf(a[0]) == objStr {
		in := strings.Join(a, " ")
		return []*Object{p.objStr(in)}, nil
	}
	objs := []*Object{}
	for _, v := range a {
		o, err := p.parseObject(v)
		if err != nil {
			return []*Object{}, err
		}
		objs = append(objs, o)
	}
	return objs, nil
}

func (p *Parser) ReadStatement(s string) (*Statement, error) {
	var (
		ctx        string
		keyword    string
		argsRaw    = []string{}
		hasKeyword bool
		stype      StatementType
		function   string
	)
	for _, c := range s {
		if p.isWhitespace(c) {
			if !hasKeyword {
				keyword = p.id(ctx)
				hasKeyword = true
				ctx = ""
			} else {
				argsRaw = append(argsRaw, strings.TrimSpace(ctx))
				ctx = ""
			}
		} else {
			ctx += string(c)
		}
	}
	if ctx != "" {
		if !hasKeyword {
			keyword = p.id(ctx)
			hasKeyword = true
		} else {
			argsRaw = append(argsRaw, strings.TrimSpace(ctx))
		}
	}
	switch keyword {
	case "set":
		stype = stSet
	case "if":
		stype = stIf
	case "else":
		stype = stElse
	case "label":
		stype = stLabel
	default:
		stype = stCall
	}
	if stype == stCall {
		function = keyword
	}
	args, err := p.parseArgs(argsRaw)
	if err != nil {
		return nil, err
	}
	return &Statement{
		Type: stype,
		Func: function,
		Args: args,
	}, nil
}
