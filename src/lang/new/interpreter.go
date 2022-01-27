package new

import (
	"mohazit/lang"
	"mohazit/lib"
)

var (
	missFunc = lib.LazyError("interpreter: unknown function %s", "nint_missfun")
	missComp = lib.LazyError("interpreter: unknown comparator %s", "nint_misscomp")
	unexTkn  = lib.LazyError("interpreter: unexpected %s token", "nint_unextkn")
	unex     = lib.LazyError("interpreter: unexpected %s", "nint_unex")
	unimpl   = lib.LazyError("interpreter: %s unimplemented", "nint_unimpl")
)

// Interpreter reads statements from it's internal Parser and exectures them
// also stores global/local variables and lables
type Interpreter struct {
	parser *Parser
	vars   map[string]*lang.Object
	labels map[string][]Statement
}

// NewInterpreter creates an empty Interpreter, which has no code to run
func NewInterpreter() *Interpreter {
	return &Interpreter{
		parser: NewParser(),
		vars:   make(map[string]*lang.Object),
		labels: make(map[string][]Statement),
	}
}

// Source gives this Interpreter some code to run
func (i *Interpreter) Source(src string) {
	i.parser.Source(src)
}

// Do runs as many statements as possible, stopping if there's a problem
// reading the next statement (first value will be false) or if there's a
// problem executing said statement (first value will be true)
func (i *Interpreter) Do() (ok bool, err error) {
	for {
		stmt, err := i.parser.Next()
		if err != nil {
			return false, err
		}
		if stmt == nil {
			return false, nil
		}
		if err = i.exec(stmt); err != nil {
			return true, i.exec(stmt)
		}
	}
}

// exec runs a singular statement, consuming more statements if necessary
func (i *Interpreter) exec(stmt *Statement) error {
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
					return unexTkn.Get(tkn.Type.String())
				}
			} else {
				switch tkn.Type {
				case tIdent, tLiteral, tSpace:
					r = append(r, tkn)
				case tOper:
					return unimpl.Get("operator chaining")
				case tLinefeed:
					break
				default:
					return unexTkn.Get(tkn.Type.String())
				}
			}
		}
		lVal, err := i.parser.Tokens2object(l)
		if err != nil {
			return err
		}
		rVal, err := i.parser.Tokens2object(r)
		if err != nil {
			return err
		}
		c, ok := lang.Comps[op.Raw]
		if !ok {
			return missComp.Get(op.Raw)
		}
		v, err := c(lVal, rVal)
		if err != nil {
			return err
		}
		if stmt.Keyword == "unless" {
			v = !v
		}
		for {
			substmt, err := i.parser.Next()
			if err != nil {
				return err
			}
			if substmt == nil {
				break
			}
			switch substmt.Keyword {
			case "else":
				v = !v
			case "end":
				return nil
			default:
				if v {
					err := i.exec(substmt)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	case "end":
		return unex.Get("end")
	default:
		f, ok := lang.Funcs[stmt.Keyword]
		if !ok {
			return missFunc.Get(stmt.Keyword)
		}
		args, err := i.parser.Args(stmt.Args)
		if err != nil {
			return err
		}
		_, err = f(args)
		return err
		// TODO(qeaml): variables, labels and other special statements
	}
}
