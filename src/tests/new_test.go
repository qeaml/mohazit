package tests

import (
	"mohazit/lang"
	"mohazit/lib"

	"testing"
)

var gt *testing.T

func TestLexer(t *testing.T) {
	lib.Load()
	gt = t
	lang.Source("var test = 123")
	expectToken(3, "var")  // ident
	expectToken(0, " ")    // space
	expectToken(3, "test") // ident
	expectToken(0, " ")    // space
	expectToken(4, "=")    // oper
	expectToken(0, " ")    // space
	expectToken(2, "123")  // literal
	expectToken(1, "\n")
}

func expectToken(tt lang.TokenType, tr string) {
	tkn := lang.NextToken()
	if tkn == nil {
		return
	}
	gt.Log(tkn.String())
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
	lang.Source("var test = 123")
	expectStatement("var", 0, 3, 0, 4, 0, 2)
}

func expectStatement(kw string, args ...lang.TokenType) {
	stmt, err := lang.NextStmt()
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
	lang.Source("file-create deez.txt\nfile-rename deez.txt \\ deez nuts.txt\nfile-delete deez nuts.txt")
	for {
		cont, err := lang.DoAll()
		if !cont {
			break
		}
		if err != nil {
			t.Fatal(err.Error())
			break
		}
	}
}

func TestCall(t *testing.T) {
	lib.Load()
	gt = t
	lang.Source("say hello\nsay world\ndata-stream blajh\ndata-write hello world\ndata-close")
	_, err := lang.DoAll()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestIf(t *testing.T) {
	lib.Load()
	gt = t
	lang.Source("unless 1 = 3\nsay aa\nsay bb\nelse\nsay dd\nend")
	_, err := lang.DoAll()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestVar(t *testing.T) {
	lib.Load()
	gt = t
	lang.Source("global i = 123\nvar j = 321\nset k=101010")
	_, err := lang.DoAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	expectGlobalVariable("i", lang.NewInt(123))
	expectGlobalVariable("j", lang.NewInt(321))
	expectGlobalVariable("k", lang.NewInt(101010))
}

func expectGlobalVariable(name string, value *lang.Object) {
	o, ok := lang.GetGlobalVar(name)
	if !ok {
		gt.Fatalf("global varialbe %s does not exist", name)
	}
	if !o.Equals(value) {
		gt.Fatalf("global variables has value %s, but %s was expected",
			o.Repr(), value.Repr())
	}
	gt.Logf("found expected variable %s with expected value %s",
		name, value.Repr())
}

func TestFunc(t *testing.T) {
	lib.Load()
	gt = t
	lang.Source("global a = [inc] 10\nvar b = [dec] 101\n set c= [dec dec dec] 9")
	_, err := lang.DoAll()
	if err != nil {
		if perr, ok := err.(*lang.ParseError); ok {
			t.Logf("%s %s", perr.Where.String(), perr.Error())
		}
		t.Fatal(err.Error())
	}

	expectGlobalVariable("a", lang.NewInt(11))
	expectGlobalVariable("b", lang.NewInt(100))
	expectGlobalVariable("c", lang.NewInt(6))
}
