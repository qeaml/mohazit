package test

import (
	"mohazit/lang"
	"testing"
)

func TestInterpreter(t *testing.T) {
	i := lang.NewInterpreter()
	err := i.RunFile("examples/test.mhzt")
	if err != nil {
		t.Fatal(err)
	}
}
