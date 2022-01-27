package tests

import (
	"mohazit/lang/new"
	"mohazit/lib"

	"testing"
)

var gt *testing.T
var gl *new.Lexer
var gp *new.Parser
var gi *new.Interpreter

func TestLexer(t *testing.T) {
	lib.Load()
	gt = t
	gl = new.NewLexer()
	gl.Source("var test = 123")
	expectToken(3, "var")  // ident
	expectToken(0, " ")    // space
	expectToken(3, "test") // ident
	expectToken(0, " ")    // space
	expectToken(4, "=")    // oper
	expectToken(0, " ")    // space
	expectToken(2, "123")  // literal
	expectToken(1, "\n")
}

func expectToken(tt new.TokenType, tr string) {
	tkn, err := gl.Next()
	if err != nil {
		gt.Fatal(err.Error())
	}
	if tkn == nil {
		return
	}
	gt.Logf("%s token: %s", tkn.Type.String(), tkn.Raw)
	gt.Logf("%d", int(tkn.Raw[0]))
	if tkn.Type != tt {
		gt.Fatalf("wrong type, got %s, want %s",
			tkn.Type.String(), tt.String())
	}
	if tkn.Raw != tr {
		gt.Fatalf("wrong raw, got %s, want %s",
			tkn.Raw, tr)
	}
}

func TestParser(t *testing.T) {
	lib.Load()
	gt = t
	gp = new.NewParser()
	gp.Source("var test = 123")
	expectStatement("var", 0, 3, 0, 4, 0, 2)
}

func expectStatement(kw string, args ...new.TokenType) {
	stmt, err := gp.Next()
	if err != nil {
		gt.Fatal(err.Error())
	}
	gt.Logf("%s statement with %d args",
		stmt.Keyword, len(stmt.Args))
	if stmt.Keyword != kw {
		gt.Fatalf("wrong keyword, got %s, want %s",
			stmt.Keyword, kw)
	}
	if len(stmt.Args) != len(args) {
		gt.Fatalf("wrong arg count, got %d, want %d",
			len(stmt.Args), len(args))
	}
	for i := 0; i < len(stmt.Args); i++ {
		if stmt.Args[i].Type != args[i] {
			gt.Fatalf("wrong arg[%d] type, want %s, got %s",
				i, stmt.Args[i].Type.String(), args[i].String())
		}
	}
	gt.Logf("all arg types match")
}

func TestInterpreter(t *testing.T) {
	lib.Load()
	gt = t
	gi = new.NewInterpreter()
	gi.Source("file-create deez.txt\nfile-rename deez.txt \\ deez nuts.txt\nfile-delete deez nuts.txt")
	for {
		cont, err := gi.Do()
		if !cont {
			break
		}
		if err != nil {
			t.Fatal(err.Error())
			break
		}
	}
}

func TestIf(t *testing.T) {
	lib.Load()
	gt = t
	gi = new.NewInterpreter()
	gi.Source("unless 1 = 3\nsay aa\nsay bb\nelse\nsay dd\nend")
	for {
		cont, err := gi.Do()
		if !cont {
			break
		}
		if err != nil {
			t.Fatal(err.Error())
			break
		}
	}
}
