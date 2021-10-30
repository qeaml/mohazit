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
	stUnless
)

type Statement struct {
	Type    StatementType
	Func    string
	Args    []*Object
	ArgsSrc string
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
		Type: ObjNil,
	}
}

func (p *Parser) objInt(i int) *Object {
	return &Object{
		Type: ObjInt,
		IntV: i,
	}
}

func (p *Parser) objBool(b bool) *Object {
	return &Object{
		Type:  ObjBool,
		BoolV: b,
	}
}

func (p *Parser) objStr(s string) *Object {
	return &Object{
		Type: ObjStr,
		StrV: s,
	}
}

func (p *Parser) typeOf(s string) ObjectType {
	s = strings.TrimSpace(s)
	if s == "" {
		return ObjNil
	}
	switch strings.ToLower(s) {
	case "nil":
		return ObjNil
	case "true", "yes", "false", "no":
		return ObjBool
	}
	if s[0] == '-' || p.isDigit(s[0]) {
		return ObjInt
	}
	return ObjStr
}

func (p *Parser) parseObject(s string, t ObjectType) (*Object, error) {
	s = strings.TrimSpace(s)
	switch t {
	case ObjNil:
		return p.objNil(), nil
	case ObjBool:
		s = strings.ToLower(s)
		if s == "true" || s == "yes" {
			return p.objBool(true), nil
		}
		if s == "false" || s == "no" {
			return p.objBool(false), nil
		}
		return p.objNil(), errors.New("invalid boolean value: " + s)
	case ObjInt:
		i, err := strconv.Atoi(s)
		if err != nil {
			return p.objNil(), err
		}
		return p.objInt(i), nil
	case ObjStr:
		return p.objStr(s), nil
	}
	return p.objNil(), errors.New("could not deterime type of value: " + s)
}

func (p *Parser) parseArgs(a []string) ([]*Object, error) {
	objs := []*Object{}
	for _, v := range a {
		t := p.typeOf(v)
		if len(objs) > 0 && t == ObjStr {
			prev := objs[len(objs)-1]
			if prev.Type == ObjStr {
				if strings.HasSuffix(prev.StrV, "\\") {
					prev.StrV = strings.TrimSpace(strings.TrimSuffix(prev.StrV, "\\"))
					objs[len(objs)-1] = prev
				} else {
					prev.StrV += " " + v
					objs[len(objs)-1] = prev
					continue
				}
			}
		}
		o, err := p.parseObject(v, t)
		if err != nil {
			return []*Object{}, err
		}
		objs = append(objs, o)
	}
	return objs, nil
}

func (p *Parser) ParseStatement(s string) (*Statement, error) {
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
	case "unless":
		stype = stUnless
	default:
		stype = stCall
	}
	if stype == stCall {
		function = keyword
	}
	args, err := p.parseArgs(argsRaw)
	if err != nil {
		return nil, p.errOf(err)
	}
	return &Statement{
		Type:    stype,
		Func:    function,
		Args:    args,
		ArgsSrc: strings.Join(argsRaw, " "),
	}, nil
}

func (p *Parser) errOf(err error) error {
	return errors.New("parser error: " + err.Error())
}
