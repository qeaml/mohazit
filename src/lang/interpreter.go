package lang

import (
	"errors"
	"fmt"
	"mohazit/tool"
	"strings"
)

type Context uint8

const (
	ctxGlobal Context = iota
	ctxIf
	ctxElse
	ctxLabel
)

type Interpreter struct {
	parser      *Parser
	ctx         Context
	ifCondition bool
	ifBlock     []*Statement
	elseBlock   []*Statement
	labelName   string
	labelBlock  []*Statement
	labelMap    map[string][]*Statement
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		parser:   &Parser{},
		labelMap: make(map[string][]*Statement),
	}
}

func (i *Interpreter) RunAll(lines []*Statement) error {
	var err error
	for _, s := range lines {
		if err = i.RunStatement(s); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) RunLine(line string) error {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil
	}
	st, err := i.parser.ReadStatement(line)
	if err != nil {
		return err
	}
	return i.RunStatement(st)
}

func (i *Interpreter) noStmt() []*Statement {
	return []*Statement{}
}

func (i *Interpreter) RunStatement(st *Statement) error {
	switch st.Type {
	case stCall:
		switch i.ctx {
		case ctxIf:
			if st.Func == "end" {
				i.ctx = ctxGlobal
				if i.ifCondition {
					tool.Log("running if block")
					i.RunAll(i.ifBlock)
					i.ifBlock = i.noStmt()
					i.elseBlock = i.noStmt()
				}
				tool.Log("if block ends")
			} else {
				i.ifBlock = append(i.ifBlock, st)
				tool.Log("if block statement: " + st.Repr())
			}
		case ctxElse:
			if st.Func == "end" {
				i.ctx = ctxGlobal
				if !i.ifCondition {
					tool.Log("running else block")
					i.RunAll(i.elseBlock)
					i.ifBlock = i.noStmt()
					i.elseBlock = i.noStmt()
				}
				tool.Log("else block ends")
			} else {
				i.elseBlock = append(i.elseBlock, st)
				tool.Log("else block statement: " + st.Repr())
			}
		case ctxLabel:
			if st.Func == "end" {
				i.ctx = ctxGlobal
				i.labelMap[i.labelName] = i.labelBlock
				i.labelBlock = i.noStmt()
				tool.Log("label block ends")
			} else {
				i.labelBlock = append(i.labelBlock, st)
				tool.Log("label block statement: " + st.Repr())
			}
		case ctxGlobal:
			tool.Log("regular function call: " + st.Repr())
			switch st.Func {
			case "goto":
				if len(st.Args) < 1 {
					return errors.New("goto requires a label name to jump to")
				}
				target := st.Args[0]
				if target.Type != objStr {
					return errors.New("goto argument must be a string")
				}
				stmt, ok := i.labelMap[target.StrV]
				if !ok {
					return errors.New("label not found: " + target.StrV)
				}
				tool.Log("GOTO going to to: " + target.StrV)
				return i.RunAll(stmt)
			case "say":
				txt := []string{}
				for _, a := range st.Args {
					txt = append(txt, a.Repr())
				}
				fmt.Println(strings.Join(txt, " "))
				return nil
			}
		}
	case stSet:
		break
	case stIf:
		if i.ctx == ctxIf {
			return errors.New("nested if blocks are not supported")
		} else if i.ctx != ctxGlobal {
			return errors.New("unexpected if block")
		}
		i.ctx = ctxIf
		tool.Log("if block starts")
	case stElse:
		if i.ctx != ctxIf {
			return errors.New("else block with no preceding if block")
		}
		i.ctx = ctxElse
		tool.Log("else block starts")
	case stLabel:
		if i.ctx != ctxGlobal {
			return errors.New("unexpected label block")
		}
		i.ctx = ctxLabel
		if len(st.Args) < 1 {
			return errors.New("must specify label name")
		}
		lnr := st.Args[0]
		if lnr.Type != objStr {
			return errors.New("label name must be a string")
		}
		i.labelName = lnr.StrV
		tool.Log("label block starts: " + i.labelName)
	}
	return nil
}
