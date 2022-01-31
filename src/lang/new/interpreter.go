package new

import (
	"fmt"
	"mohazit/lang"
)

var globals = make(map[string]*lang.Object)
var locals = make(map[string]*lang.Object)
var labels = make(map[string][]Statement)

// DoAll runs as many statements as possible, stopping if there's a problem
// reading the next statement (first value will be false) or if there's a
// problem executing said statement (first value will be true)
func DoAll() (ok bool, err error) {
	for {
		stmt, err := NextStmt()
		if err != nil {
			return false, err
		}
		if stmt == nil {
			return false, nil
		}
		if err = RunStmt(stmt, false); err != nil {
			return true, err
		}
	}
}

// RunStmt runs a singular statement, consuming more statements if necessary
func RunStmt(stmt *Statement, isLocal bool) error {
	switch stmt.Keyword {
	case "if", "unless":
		if !isLocal { // don't naively wipe locals
			locals = make(map[string]*lang.Object)
		}
		l := []*Token{}
		var op *Token = nil
		r := []*Token{}
	ifLoop:
		for _, tkn := range stmt.Args {
			if op == nil {
				switch tkn.Type {
				case tIdent, tLiteral, tSpace:
					l = append(l, tkn)
				case tOper:
					op = tkn
				case tLinefeed:
					break ifLoop
				default:
					return perrf(tkn, "unexpected %s", tkn.Type.String())
				}
			} else {
				switch tkn.Type {
				case tIdent, tLiteral, tSpace:
					r = append(r, tkn)
				case tOper:
					return perr(tkn, "operator chaining not yet implemented")
				case tLinefeed:
					break ifLoop
				default:
					return perrf(tkn, "unexpected %s", tkn.Type.String())
				}
			}
		}
		lVal, err := parseObject(l)
		if err != nil {
			return err
		}
		rVal, err := parseObject(r)
		if err != nil {
			return err
		}
		c, ok := lang.Comps[op.Raw]
		if !ok {
			return perrf(op, "unknown comparator %s", op.Raw)
		}
		v, err := c(lVal, rVal)
		if err != nil {
			return err
		}
		if stmt.Keyword == "unless" {
			v = !v
		}
		for {
			substmt, err := NextStmt()
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
					err := RunStmt(substmt, true)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	case "end":
		return perr(stmt.KwToken, "end statement outside of block")
	case "local", "global", "var", "set":
		l := []*Token{}
		mid := false
		r := []*Token{}
	varLoop:
		for _, tkn := range stmt.Args {
			if !mid {
				switch tkn.Type {
				case tIdent, tLiteral, tSpace:
					l = append(l, tkn)
				case tOper:
					if tkn.Raw == "=" {
						mid = true
					} else {
						return perrf(tkn, "expected =, got %s", tkn.Raw)
					}
				case tLinefeed:
					break varLoop
				default:
					return perrf(tkn, "unexpected %s", tkn.Type.String())
				}
			} else {
				switch tkn.Type {
				case tIdent, tLiteral, tSpace:
					r = append(r, tkn)
				case tLinefeed:
					break varLoop
				default:
					return perrf(tkn, "unexpected %s", tkn.Type.String())
				}
			}
		}
		l = trimSpaceTokens(l)
		if len(l) > 1 {
			return perr(l[0], "too many tokens before =")
		}
		lVal := l[0]
		if lVal.Type != tIdent {
			fmt.Println(l)
			return perrf(lVal, "expected identifier, got %s", lVal.Type.String())
		}
		rVal, err := parseObject(r)
		if err != nil {
			return err
		}
		if stmt.Keyword == "local" {
			if !isLocal {
				return perr(stmt.KwToken, "local variable in global context")
			}
			locals[lVal.Raw] = rVal
		} else if stmt.Keyword == "global" {
			globals[lVal.Raw] = rVal
		} else {
			if isLocal {
				locals[lVal.Raw] = rVal
			} else {
				globals[lVal.Raw] = rVal
			}
		}
		return nil
	default:
		f, ok := lang.Funcs[stmt.Keyword]
		if !ok {
			return perrf(stmt.KwToken, "unknown function %s", stmt.Keyword)
		}
		args, err := parseObjectList(stmt.Args)
		if err != nil {
			return err
		}
		_, err = f(args)
		return err
		// TODO(qeaml): variables, labels and other special statements
	}
}

func GetGlobalVar(name string) (v *lang.Object, ok bool) {
	v, ok = globals[name]
	return
}
