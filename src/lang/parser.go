package lang

import (
	"errors"
	"mohazit/tool"
	"strconv"
	"strings"
)

func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

type parser struct {
}

func (p *parser) isWhitespace(c rune) bool {
	return c == ' ' || c == '\t'
}

func (p *parser) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (p *parser) objNil() *Object {
	return &Object{
		Type: ObjNil,
	}
}

func (p *parser) objInt(i int) *Object {
	return &Object{
		Type: ObjInt,
		IntV: i,
	}
}

func (p *parser) objBool(b bool) *Object {
	return &Object{
		Type:  ObjBool,
		BoolV: b,
	}
}

func (p *parser) objStr(s string) *Object {
	return &Object{
		Type: ObjStr,
		StrV: s,
	}
}

func (p *parser) typeOf(s string) ObjectType {
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

func (p *parser) parseObject(s string, t ObjectType) (*Object, error) {
	assert(t < ObjInv, "object type invalid")
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

func (p *parser) parseArgs(a []string) ([]*Object, error) {
	tool.Log("parse args call: ", a)
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

type genStmt struct {
	Kw  string
	Arg string
}

func (p *parser) ParseStatement(s string) (*genStmt, error) {
	ctx := ""
	hasKw := false
	kw := ""
	// main parsing loop: for each character,
	for _, c := range s {
		// if we don't have the keyword and the current char is whitespace
		if !hasKw && p.isWhitespace(c) {
			// then everything we've read so far is the keyword
			kw = strings.ToLower(strings.TrimSpace(ctx))
			hasKw = true
			ctx = ""
		} else {
			// otherwise just add it to the context
			ctx += string(c)
		}
	}
	// keyword-only statement (else, end etc.)
	if !hasKw {
		kw = strings.ToLower(strings.TrimSpace(ctx))
		hasKw = true
		ctx = ""
	}
	tool.Log("- Out: " + kw + "(" + ctx + ")")
	return &genStmt{
		Kw:  kw,
		Arg: strings.TrimSpace(ctx),
	}, nil
}

type condStmt struct {
	Kw   string
	Cond *Conditional
}

func (p *parser) toCond(gs *genStmt) (*condStmt, error) {
	condition, err := ParseConditional(gs.Arg, p)
	if err != nil {
		return nil, err
	}
	return &condStmt{
		Kw:   gs.Kw,
		Cond: condition,
	}, nil
}

type callStmt struct {
	Kw   string
	Args []*Object
}

func (p *parser) toCall(gs *genStmt) (*callStmt, error) {
	args, err := p.parseArgs([]string{gs.Arg})
	if err != nil {
		return nil, err
	}
	return &callStmt{
		Kw:   gs.Kw,
		Args: args,
	}, nil
}
