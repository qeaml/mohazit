package new

import (
	"mohazit/lang"
	"mohazit/lib"
	"strconv"
)

var (
	missFunc = lib.LazyError("interpreter: unknown function %s", "nint_missfun")
	missComp = lib.LazyError("interpreter: unknown comparator %s", "nint_misscomp")
	tooMany  = lib.LazyError("interpreter: too many Do(%s)s", "nint_toomany")
	unexTkn  = lib.LazyError("interpreter: unexpected %s token", "nint_unextkn")
	unex     = lib.LazyError("interpreter: unexpected %s", "nint_unex")
	unimpl   = lib.LazyError("interpreter: %s unimplemented", "nint_unimpl")
)

type Interpreter struct {
	parser *Parser
	vars   map[string]*lang.Object
	labels map[string][]Statement
}

func NewInterpreter(source string) *Interpreter {
	lex := NewLexer(source)
	return &Interpreter{
		parser: &Parser{lex},
		vars:   make(map[string]*lang.Object),
		labels: make(map[string][]Statement),
	}
}

func (i *Interpreter) Do() (bool, error) {
	stmt, err := i.parser.Next()
	if err != nil {
		return true, err
	}
	if stmt == nil {
		return false, tooMany.Get("")
	}
	return i.exec(stmt)
}

func (i *Interpreter) exec(stmt *Statement) (bool, error) {
	switch stmt.Keyword {
	case "if", "unless":
		l := []*Token{}
		var op *Token = nil
		r := []*Token{}
		for _, tkn := range stmt.Args {
			if op == nil {
				switch tkn.Type {
				case tIdent, tLiteral, tSpace:
					l = append(l, tkn)
				case tOper:
					op = tkn
				case tLinefeed:
					break
				default:
					return true, unexTkn.Get(tkn.Type.String())
				}
			} else {
				switch tkn.Type {
				case tIdent, tLiteral, tSpace:
					r = append(r, tkn)
				case tOper:
					return true, unimpl.Get("operator chaining")
				case tLinefeed:
					break
				default:
					return true, unexTkn.Get(tkn.Type.String())
				}
			}
		}
		lVal, err := i.tokens2object(l)
		if err != nil {
			return true, err
		}
		rVal, err := i.tokens2object(r)
		if err != nil {
			return true, err
		}
		c, ok := lang.Comps[op.Raw]
		if !ok {
			return true, missComp.Get(op.Raw)
		}
		v, err := c(lVal, rVal)
		if err != nil {
			return true, err
		}
		if stmt.Keyword == "unless" {
			v = !v
		}
		for {
			substmt, err := i.parser.Next()
			if err != nil {
				return true, err
			}
			if substmt == nil {
				break
			}
			switch substmt.Keyword {
			case "else":
				v = !v
			case "end":
				return true, nil
			default:
				if v {
					_, err := i.exec(substmt)
					if err != nil {
						return true, err
					}
				}
			}
		}
		return true, nil
	case "end":
		return true, unex.Get("end")
	default:
		f, ok := lang.Funcs[stmt.Keyword]
		if !ok {
			return true, missFunc.Get(stmt.Keyword)
		}
		args, err := i.parser.Args(stmt.Args)
		if err != nil {
			return true, err
		}
		_, err = f(args)
		return true, err
		// TODO(qeaml): variables, labels and other special statements
	}
}

func (i *Interpreter) tokens2object(t []*Token) (*lang.Object, error) {
	t = i.trimSpace(t)
	switch t[0].Type {
	case tIdent, tInvalid:
		v := lang.NewStr(t[0].Raw)
		for i := 0; i < len(t); i++ {
			tkn := t[i]
			switch tkn.Type {
			case tIdent, tInvalid, tSpace:
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

func (i *Interpreter) trimSpace(t []*Token) []*Token {
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
