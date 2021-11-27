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
}

func NewInterpreter(p *parser) *interpreter {
	return &interpreter{
		parser:   p,
		labelMap: make(map[string][]*genStmt),
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
	// TODO(qeaml): reference processing
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

func (i *interpreter) runGlobally(st *genStmt) error {
	switch st.Kw {
	case "if":
		condSt, err := i.parser.toCond(st)
		if err != nil {
			return err
		}
		comp, ok := Comps[condSt.Cond.Comp]
		if !ok {
			return i.err("unknown comparator: " + condSt.Cond.Comp)
		}
		i.condBlock = i.noStmt()
		i.cond, err = comp(condSt.Cond.Args)
		if err != nil {
			return i.errOf(err)
		}
		i.ctx = ctxIf
		return nil
	case "unless":
		condSt, err := i.parser.toCond(st)
		if err != nil {
			return err
		}
		comp, ok := Comps[condSt.Cond.Comp]
		if !ok {
			return i.err("unknown comparator: " + condSt.Cond.Comp)
		}
		i.condBlock = i.noStmt()
		i.cond, err = comp(condSt.Cond.Args)
		if err != nil {
			return i.errOf(err)
		}
		i.cond = !i.cond
		i.ctx = ctxUnless
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
		condSt, err := i.parser.toCond(st)
		if err != nil {
			return err
		}
		comp, ok := Comps[condSt.Cond.Comp]
		if !ok {
			return i.err("unknown comparator: " + condSt.Cond.Comp)
		}
		res, err := comp(condSt.Cond.Args)
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
		condSt, err := i.parser.toCond(st)
		if err != nil {
			return err
		}
		comp, ok := Comps[condSt.Cond.Comp]
		if !ok {
			return i.err("unknown comparator: " + condSt.Cond.Comp)
		}
		res, err := comp(condSt.Cond.Args)
		if err != nil {
			return err
		}
		if res {
			fmt.Println("FAIL [Not] " + st.Arg)
			return i.err("assertion didn't fail")
		}
		fmt.Println("PASS [Not] " + st.Arg)
		return nil
	default:
		callSt, err := i.parser.toCall(st)
		if err != nil {
			return err
		}
		f, ok := Funcs[callSt.Kw]
		if !ok {
			return i.err("unknown function: " + callSt.Kw)
		}
		return f(callSt.Args)
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

type FuncMap map[string]func([]*Object) error
type CompMap map[string]func([]*Object) (bool, error)

var Funcs = make(FuncMap)
var Comps = make(CompMap)
