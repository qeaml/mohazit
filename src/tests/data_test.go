package tests

import (
	"mohazit/lang"
	"mohazit/lib"
	"testing"
)

func TestBuf(t *testing.T) {
	gt = t
	lib.Load()
	lang.Source(`
		var b = [buf-create]
		data-write hello
		data-seek 1
		var res = [data-read] 4
	`)
	err := lang.DoAll()
	if err != nil {
		if perr, ok := err.(*lang.ParseError); ok {
			t.Fatalf("%s @%s", perr.Error(), perr.Where)
		} else {
			t.Fatal(err.Error())
		}
	}

	expectGlobalVariable("res", "ello")
}

func TestFile(t *testing.T) {
	gt = t
	lib.Load()
	lang.Source(`
		file-create test.txt
		var f = [file-open] test.txt
		data-write he world
		data-seek 2
		data-write llo
		data-close
		var f = [file-open] test.txt
		data-seek 1
		var res = [data-read] 4
		data-close
		file-delete test.txt
	`)
	err := lang.DoAll()
	if err != nil {
		if perr, ok := err.(*lang.ParseError); ok {
			t.Fatalf("%s @%s", perr.Error(), perr.Where)
		} else {
			t.Fatal(err.Error())
		}
	}

	expectGlobalVariable("res", "ello")
}
