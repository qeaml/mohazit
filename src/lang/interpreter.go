package lang

import (
	"errors"
	"fmt"
)

type Context uint8

const (
	ctxGlobal Context = iota
	ctxIf
	ctxElse
	ctxLabel
	ctxUnless
)

type interpreter struct {
	parser     *parser
	ctx        Context
	cond       bool
	condBlock  []*genStmt
	elseBlock  []*genStmt
	labelName  string
	labelBlock []*genStmt
	labelMap   map[string][]*genStmt
	vars       map[string]*Object
}

func NewInterpreter(p *parser) *interpreter {
	return &interpreter{
		parser:   p,
		labelMap: make(map[string][]*genStmt),
		vars:     make(map[string]*Object),
	}
}

func (i *interpreter) runAll(st []*genStmt) error {
	var err error
	for _, s := range st {
		if err = i.RunStatement(s); err != nil {
			return err
		}
	}
	return nil
}

func (i *interpreter) noStmt() []*genStmt {
	return []*genStmt{}
}

func (i *interpreter) RunStatement(st *genStmt) error {
	switch i.ctx {
	case ctxGlobal:
		return i.runGlobally(st)
	case ctxIf:
		return i.runInIf(st)
	case ctxElse:
		return i.runInElse(st)
	case ctxUnless:
		return i.runInUnless(st)
	case ctxLabel:
		return i.runInLabel(st)
	default:
		return i.err("internal interpreter error")
	}
}

type InterVar struct {
	Set   func(string, *Object)
	Get   func(string) (*Object, bool)
	Has   func(string) bool
	Parse func(string) (*Object, error)
}

func (i *interpreter) runGlobally(st *genStmt) error {
	switch st.Kw {
	case "if", "unless":
		condSt, err := i.parser.toCond(st, i.vars)
		if err != nil {
			return err
		}
		comp, ok := Comps[condSt.Comp]
		if !ok {
			return i.err("unknown comparator: " + condSt.Comp)
		}
		i.condBlock = i.noStmt()
		i.cond, err = comp(condSt.Args)
		if err != nil {
			return i.errOf(err)
		}
		if st.Kw == "if" {
			i.ctx = ctxIf
		} else if st.Kw == "unless" {
			i.ctx = ctxUnless
			i.cond = !i.cond
		}
		return nil
	case "label":
		i.labelName = st.Arg
		i.ctx = ctxLabel
		return nil
	case "else":
		return i.err("else outside if/unless block")
	case "end":
		return i.err("end outside block")
	case "goto":
		if len(st.Arg) < 1 {
			return i.err("must prove label name to go to")
		}
		label, ok := i.labelMap[st.Arg]
		if !ok {
			return i.err("unknown label: " + st.Arg)
		}
		return i.runAll(label)
	case "assert":
		condSt, err := i.parser.toCond(st, i.vars)
		if err != nil {
			return err
		}
		comp, ok := Comps[condSt.Comp]
		if !ok {
			return i.err("unknown comparator: " + condSt.Comp)
		}
		res, err := comp(condSt.Args)
		if err != nil {
			return err
		}
		if !res {
			fmt.Println("FAIL " + st.Arg)
			return i.err("assertion failed")
		}
		fmt.Println("PASS " + st.Arg)
		return nil
	case "assert-not":
		condSt, err := i.parser.toCond(st, i.vars)
		if err != nil {
			return err
		}
		comp, ok := Comps[condSt.Comp]
		if !ok {
			return i.err("unknown comparator: " + condSt.Comp)
		}
		res, err := comp(condSt.Args)
		if err != nil {
			return err
		}
		if res {
			fmt.Println("FAIL [Not] " + st.Arg)
			return i.err("assertion didn't fail")
		}
		fmt.Println("PASS [Not] " + st.Arg)
		return nil
	case "set", "var":
		varSt, err := i.parser.toVar(st)
		if err != nil {
			return err
		}
		if !varSt.Processed {
			i.vars[varSt.Name] = varSt.Value
		} else {
			proc, ok := Procs[varSt.Processor]
			if !ok {
				return i.err("unknown processor: " + varSt.Processor)
			}
			res, err := proc(varSt.Value)
			if err != nil {
				return err
			}
			i.vars[varSt.Name] = res
		}
		return nil
	case "dump-var":
		fmt.Println(st.Arg, "=", i.vars[st.Arg].Repr())
		return nil
	default:
		callSt, err := i.parser.toCall(st, i.vars)
		if err != nil {
			return err
		}
		f, ok := Funcs[callSt.Kw]
		if !ok {
			return i.err("unknown function: " + callSt.Kw)
		}
		return f(callSt.Args, InterVar{
			Get: func(s string) (*Object, bool) {
				o, ok := i.vars[s]
				return o, ok
			},
			Set: func(s string, o *Object) {
				i.vars[s] = o
			},
			Has: func(s string) bool {
				_, ok := i.vars[s]
				return ok
			},
			Parse: func(s string) (*Object, error) {
				o, err := i.parser.parseObject(s, i.parser.typeOf(s))
				return o, err
			}})
	}
}

func (i *interpreter) runInIf(st *genStmt) error {
	switch st.Kw {
	case "else":
		i.elseBlock = i.noStmt()
		i.ctx = ctxElse
		return nil
	case "end":
		i.ctx = ctxGlobal
		if i.cond {
			return i.runAll(i.condBlock)
		}
	case "if":
		return i.err("nested if blocks are not supported")
	case "unless":
		return i.err("nested unless blocks are not supported")
	case "label":
		return i.err("labels cannot be defined conditionally")
	default:
		i.condBlock = append(i.condBlock, st)
	}
	return nil
}

func (i *interpreter) runInElse(st *genStmt) error {
	switch st.Kw {
	case "else":
		return i.err("nested else blocks are not supported")
	case "end":
		i.ctx = ctxGlobal
		if !i.cond {
			return i.runAll(i.elseBlock)
		} else {
			return i.runAll(i.condBlock)
		}
	case "if":
		return i.err("nested if blocks are not supported")
	case "unless":
		return i.err("nested unless blocks are not supported")
	case "label":
		return i.err("labels cannot be defined conditionally")
	default:
		i.elseBlock = append(i.elseBlock, st)
	}
	return nil
}

func (i *interpreter) runInUnless(st *genStmt) error {
	switch st.Kw {
	case "else":
		i.ctx = ctxElse
		i.elseBlock = i.noStmt()
	case "end":
		i.ctx = ctxGlobal
		if i.cond {
			return i.runAll(i.condBlock)
		}
	case "if":
		return i.err("nested if blocks are not supported")
	case "unless":
		return i.err("nested unless blocks are not supported")
	case "label":
		return i.err("labels cannot be defined conditionally")
	default:
		i.condBlock = append(i.condBlock, st)
	}
	return nil
}

func (i *interpreter) runInLabel(st *genStmt) error {
	switch st.Kw {
	case "end":
		i.ctx = ctxGlobal
		i.labelMap[i.labelName] = i.labelBlock
	default:
		i.labelBlock = append(i.labelBlock, st)
	}
	return nil
}

func (i *interpreter) errOf(err error) error {
	return i.err("interpreter error: " + err.Error())
}

func (i *interpreter) err(txt string) error {
	return errors.New("interpreter error: " + txt)
}

type FuncMap map[string]func([]*Object, InterVar) error
type CompMap map[string]func([]*Object) (bool, error)
type ProcMap map[string]func(*Object) (*Object, error)

var Funcs = make(FuncMap)
var Comps = make(CompMap)
var Procs = make(ProcMap)
