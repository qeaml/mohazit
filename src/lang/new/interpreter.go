package new

import (
	"fmt"
	"mohazit/lang"
	"mohazit/lib"
	"strings"
)

var (
	unknownFunc = lib.LazyError("interpreter: unknown function %s", "nint_unknownfunc")
	needMore    = lib.LazyError("interpreter: one too many Do(%s)s", "nint_needmore")
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
		return false, needMore.Get("")
	}
	switch stmt.Keyword {
	default:
		f, ok := lang.Funcs[stmt.Keyword]
		if !ok {
			return true, unknownFunc.Get(stmt.Keyword)
		}
		args, err := i.parser.Args(stmt.Args)
		if err != nil {
			return true, err
		}
		argsStrings := []string{}
		for _, o := range args {
			argsStrings = append(argsStrings, o.String())
		}
		fmt.Printf("%s(%s)\n", stmt.Keyword, strings.Join(argsStrings, ", "))
		_, err = f(args)
		return true, err
		// TODO(qeaml): control flow, variables and everything else
	}
}
