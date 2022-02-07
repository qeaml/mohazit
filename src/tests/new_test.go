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
	lang.Source("var test = 123\n")

	expectToken(3, "var")  // ident
	expectToken(0, " ")    // space
	expectToken(3, "test") // ident
	expectToken(0, " ")    // space
	expectToken(4, "=")    // oper
	expectToken(0, " ")    // space
	expectToken(2, "123")  // literal
	expectToken(1, "\n")

	lang.Source("say {deez} nuts")

	expectToken(3, "say")
	expectToken(0, " ")
	expectToken(6, "deez")
	expectToken(0, " ")
	expectToken(3, "nuts")
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
	lang.Source(`
		file-create deez.txt
		file-rename deez.txt \ deez nuts.txt
		file-delete deez nuts.txt
	`)
	err := lang.DoAll()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCall(t *testing.T) {
	lib.Load()
	gt = t
	lang.Source(`
		say hello
		say world
		buf-create blajh
		data-write hello world
		data-close
	`)
	err := lang.DoAll()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestIf(t *testing.T) {
	lib.Load()
	gt = t
	lang.Source(`
		if 3 = 3
			global a-ok = true
		end
		if 1 = 3
			say whoa
		else
			global b-ok = true
		end
		unless 1 = 3
			global c-ok = true
		end
		var my-var = 100000
		if 1 > {my-var}
			say woha!!!
		else
			global d-ok = true
		end
	`)
	err := lang.DoAll()
	if err != nil {
		if perr, ok := err.(*lang.ParseError); ok {
			t.Fatalf("%s @%s", perr.Error(), perr.Where)
		} else {
			t.Fatal(err.Error())
		}
	}

	expectGlobalVariable("a-ok", true)
	expectGlobalVariable("b-ok", true)
	expectGlobalVariable("c-ok", true)
	expectGlobalVariable("d-ok", true)
}

func TestVar(t *testing.T) {
	lib.Load()
	gt = t
	lang.Source(`
		global i = 123
		var j = 321
		set k=101010
	`)
	err := lang.DoAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	expectGlobalVariable("i", 123)
	expectGlobalVariable("j", 321)
	expectGlobalVariable("k", 101010)

	lang.Source("set l = {i}")
	err = lang.DoAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	expectGlobalVariable("l", 123)
}

func expectGlobalVariable(name string, value interface{}) {
	obj := lang.NewObject(value)
	o, ok := lang.GetGlobalVar(name)
	if !ok {
		gt.Fatalf("global varialbe %s does not exist", name)
	}
	if o == nil {
		gt.Fatal("variable is stored as nil")
	}
	if !o.Equals(obj) {
		gt.Fatalf("global variable %s has value %s, but %s was expected",
			name, o.Repr(), obj.Repr())
	}
	gt.Logf("found expected variable %s with expected value %s",
		name, obj.Repr())
}

func TestFunc(t *testing.T) {
	lib.Load()
	gt = t
	lang.Source(`
		global a = [inc] 10
		var b= [dec] 101
		set c=[dec dec dec] 9
	`)
	err := lang.DoAll()
	if err != nil {
		if perr, ok := err.(*lang.ParseError); ok {
			t.Logf("%s %s", perr.Where.String(), perr.Error())
		}
		t.Fatal(err.Error())
	}

	expectGlobalVariable("a", 11)
	expectGlobalVariable("b", 100)
	expectGlobalVariable("c", 6)
}

func TestObject(t *testing.T) {
	lib.Load()
	gt = t
	lang.Source(`
		global n = nil
		global i = 123
		global s = hello
		global b = true
	`)
	err := lang.DoAll()
	if err != nil {
		if perr, ok := err.(*lang.ParseError); ok {
			t.Logf("%s %s", perr.Where.String(), perr.Error())
		}
		t.Fatal(err.Error())
	}

	expectGlobalVariable("n", nil)
	expectGlobalVariable("i", 123)
	expectGlobalVariable("s", "hello")
	expectGlobalVariable("b", true)
}

func TestLabel(t *testing.T) {
	lib.Load()
	gt = t
	lang.Source(`
		label hello-world
			global ok = true
		end
		goto hello-world
	`)
	err := lang.DoAll()
	if err != nil {
		if perr, ok := err.(*lang.ParseError); ok {
			t.Logf("%s %s", perr.Where.String(), perr.Error())
		}
		t.Fatal(err.Error())
	}

	expectGlobalVariable("ok", true)
}
