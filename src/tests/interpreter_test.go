package test

import (
	"mohazit/lang"
	"testing"
)

func TestInterpreter(t *testing.T) {
	r := lang.NewRunner("test.mhzt")
	err := r.DoFile("examples/test.mhzt")
	if err != nil {
		t.Fatal(err)
	}
}
