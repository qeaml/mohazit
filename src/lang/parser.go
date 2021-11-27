package lang

import (
	"errors"
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
	if len(s) == 0 {
		panic("invalid value!")
	}
	if strings.HasPrefix(s, "\\(") {
		return ObjRef
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

func (p *parser) parseArgs(a string) ([]*Object, error) {
	// TODO(qeaml): reference parsing
	out := []*Object{}
	if a == "" {
		return out, nil
	}
	ctx := ""
	a += " "
	var obj *Object
	for _, c := range a {
		if p.isWhitespace(c) {
			v := strings.TrimSpace(ctx)
			if len(v) == 0 {
				continue
			}
			t := p.typeOf(v)
			switch t {
			case ObjStr:
				if len(out) >= 1 && !strings.HasSuffix(v, "\\") {
					obj = out[len(out)-1]
					if obj.Type == ObjStr {
						obj.StrV = strings.TrimSpace(obj.StrV + " " + v)
						out[len(out)-1] = obj
					} else {
						obj = p.objStr(v)
						out = append(out, obj)
					}
				} else {
					obj = p.objStr(strings.TrimSpace(strings.TrimSuffix(v, "\\")))
					out = append(out, obj)
				}
			default:
				obj, err := p.parseObject(v, t)
				if err != nil {
					return nil, err
				}
				out = append(out, obj)
			}
			ctx = ""
		} else {
			ctx += string(c)
		}
	}
	return out, nil
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
	return &genStmt{
		Kw:  kw,
		Arg: strings.TrimSpace(ctx),
	}, nil
}

type condStmt struct {
	Kw   string
	Comp string
	Args []*Object
}

func (p *parser) toCond(gs *genStmt, vars map[string]*Object) (*condStmt, error) {
	comp, args, err := p.parseConditional(gs.Arg)
	if err != nil {
		return nil, err
	}
	finalArgs := []*Object{}
	for _, a := range args {
		if a.Type == ObjRef {
			val, ok := vars[a.RefV]
			if !ok {
				return nil, errors.New("unknown variable: " + a.RefV)
			}
			finalArgs = append(finalArgs, val)
		} else {
			finalArgs = append(finalArgs, a)
		}
	}
	return &condStmt{
		Kw:   gs.Kw,
		Comp: comp,
		Args: finalArgs,
	}, nil
}

func (p *parser) parseConditional(s string) (string, []*Object, error) {
	ctx := ""
	params := []string{}
	hasParams := false
	comp := ""
	hasComp := false
	for _, c := range s {
		if hasComp || !hasParams {
			if c == '(' {
				if len(params) >= 1 && (p.typeOf(params[len(params)-1]) == ObjStr) {
					params = append(params, "\\")
				}
				hasParams = true
			} else if c == ' ' {
				a := strings.TrimSpace(ctx)
				if len(a) == 0 {
					continue
				}
				params = append(params, a)
				ctx = ""
			} else {
				ctx += string(c)
			}
		} else {
			if c == ')' {
				comp = strings.ToLower(strings.TrimSpace(comp))
				hasComp = true
			} else {
				comp += string(c)
			}
		}
	}
	if len(strings.TrimSpace(ctx)) != 0 && hasComp {
		params = append(params, strings.TrimSpace(ctx))
	}
	if !hasComp {
		return "", nil, errors.New("no comparator specified")
	}
	args, err := p.parseArgs(strings.Join(params, " "))
	if err != nil {
		return "", nil, err
	}
	return comp, args, nil
}

type callStmt struct {
	Kw   string
	Args []*Object
}

func (p *parser) toCall(gs *genStmt, vars map[string]*Object) (*callStmt, error) {
	args, err := p.parseArgs(gs.Arg)
	if err != nil {
		return nil, err
	}
	finalArgs := []*Object{}
	for _, a := range args {
		if a.Type == ObjRef {
			val, ok := vars[a.RefV]
			if !ok {
				return nil, errors.New("unknown variable: " + a.RefV)
			}
			finalArgs = append(finalArgs, val)
		} else {
			finalArgs = append(finalArgs, a)
		}
	}
	return &callStmt{
		Kw:   gs.Kw,
		Args: finalArgs,
	}, nil
}

type varStmt struct {
	name  string
	value *Object
}

func (p *parser) toVar(gs *genStmt) (*varStmt, error) {
	name := ""
	hasName := false
	valueRaw := ""
	for _, c := range gs.Arg {
		if !hasName {
			if c == '=' {
				hasName = true
				continue
			}
			name += string(c)
		} else {
			valueRaw += string(c)
		}
	}
	if !hasName {
		return nil, errors.New("variables must have a value")
	}
	values, err := p.parseArgs(valueRaw)
	if err != nil {
		return nil, err
	}
	if len(values) > 1 {
		return nil, errors.New("variables can only have 1 value")
	}
	return &varStmt{
		name:  strings.ToLower(strings.TrimSpace(name)),
		value: values[0],
	}, nil
}
