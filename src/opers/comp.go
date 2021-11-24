package opers

import (
	"errors"
	"fmt"
	"math/rand"
	"mohazit/lang"
	"mohazit/tool"
	"time"
)

func cEquals(objs []*lang.Object) (bool, error) {
	for _, o := range objs {
		tool.Log("cEquals - Object - " + o.String())
	}
	for _, a := range objs {
		for _, b := range objs {
			tool.Log("cEquals - Comparing -", a.Repr(), b.Repr())
			if a.Type != b.Type {
				return false, nil
			}
			switch a.Type {
			case lang.ObjBool:
				if a.BoolV != b.BoolV {
					return false, nil
				}
			case lang.ObjInt:
				if a.IntV != b.IntV {
					return false, nil
				}
			case lang.ObjStr:
				if a.StrV != b.StrV {
					return false, nil
				}
			}
		}
	}
	return true, nil
}

func cNotEquals(objs []*lang.Object) (bool, error) {
	eq, err := cEquals(objs)
	return !eq, err
}

func cChance(objs []*lang.Object) (bool, error) {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	return r.Float64() < 0.5, nil
}

func fCompare(args []*lang.Object) error {
	if len(args) < 1 {
		return needArgs("comparator missing")
	}
	compName := args[0]
	if compName.Type != lang.ObjStr {
		return badType("comparator must be a string")
	}
	f, ok := lang.Comps[compName.StrV]
	if !ok {
		return errors.New("unknown comparator: " + compName.StrV)
	}
	res, err := f(args[1:])
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
