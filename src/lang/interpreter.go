package lang

import (
	"errors"
	"io"
	"mohazit/tool"
	"os"
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
	st, err := i.parser.ParseStatement(line)
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
					return i.err("goto requires a label name to jump to")
				}
				target := st.Args[0]
				if target.Type != ObjStr {
					return i.err("goto argument must be a string")
				}
				stmt, ok := i.labelMap[target.StrV]
				if !ok {
					return i.err("label not found: " + target.StrV)
				}
				tool.Log("GOTO going to to: " + target.StrV)
				return i.RunAll(stmt)
			default:
				f, ok := Funcs[st.Func]
				if !ok {
					return i.err("unknown function: " + st.Func)
				}
				return f(st.Args)
			}
		}
	case stSet:
		break
	case stIf:
		if i.ctx == ctxIf {
			return i.err("nested if blocks are not supported")
		} else if i.ctx != ctxGlobal {
			return i.err("unexpected if block")
		}
		i.ctx = ctxIf
		tool.Log("if block starts")
	case stElse:
		if i.ctx != ctxIf {
			return i.err("else block with no preceding if block")
		}
		i.ctx = ctxElse
		tool.Log("else block starts")
	case stLabel:
		if i.ctx != ctxGlobal {
			return i.err("unexpected label block")
		}
		i.ctx = ctxLabel
		if len(st.Args) < 1 {
			return i.err("must specify label name")
		}
		lnr := st.Args[0]
		if lnr.Type != ObjStr {
			return i.err("label name must be a string")
		}
		i.labelName = lnr.StrV
		tool.Log("label block starts: " + i.labelName)
	}
	return nil
}

func (i *Interpreter) RunFile(fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	srcRaw, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	src := string(srcRaw)
	lines := []string{}
	ctx := ""
	for _, c := range src {
		if c == '\n' {
			if ctx[len(ctx)-1] == '\\' {
				ctx = ctx[:len(ctx)-1]
			} else {
				lines = append(lines, ctx)
				ctx = ""
			}
		} else {
			ctx += string(c)
		}
	}
	for _, l := range lines {
		if err = i.RunLine(l); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Context() string {
	switch i.ctx {
	case ctxIf:
		return "if   "
	case ctxElse:
		return "else "
	case ctxLabel:
		n := i.labelName
		if len(n) > 4 {
			n = n[:4]
		}
		if len(n) < 4 {
			for len(n) != 4 {
				n += " "
			}
		}
		return n + " "
	}
	return ""
}

// func (i *Interpreter) errOf(err error) error {
// 	return i.err("interpreter error: " + err.Error())
// }

func (i *Interpreter) err(txt string) error {
	return errors.New("interpreter error: " + txt)
}

type FuncMap map[string]func([]*Object) error

var Funcs = make(FuncMap)
