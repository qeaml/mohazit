package test

import (
	"mohazit/lang"
	"testing"
)

var (
	itl = []string{
		"say hello",
		"say what is up",
		"say 123 -123",
		"say no yes no",
		"label win",
		"	say Congratulations!",
		"	say You win",
		"end",
		"label lose",
		"	say Oh well this is sad...",
		"   say You lost",
		"end",
		"if condition",
		"	say that is true",
		"	say how about you win now",
		"	goto win",
		"else",
		"	say that is not true",
		"	say so, you sadly do not win",
		"	goto lose",
		"end",
	}
)

func TestInterpreter(t *testing.T) {
	i := lang.NewInterpreter()
	var err error
	for n, l := range itl {
		if err = i.RunLine(l); err != nil {
			t.Fatalf("error on line %d: %s", n, err.Error())
		}
	}
}
