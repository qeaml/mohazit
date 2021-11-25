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
	tool.Log("cEquals - Result - true")
	return true, nil
}

func cNotEquals(objs []*lang.Object) (bool, error) {
	eq, err := cEquals(objs)
	return !eq, err
}

func cLike(objs []*lang.Object) (bool, error) {
	if len(objs) < 2 {
		return false, needArgs("need at least 2 arguments to compare")
	}
	// target type that we will try to convert to
	tt := objs[0].Type
	// converted objects we will pass to cEquals
	co := []*lang.Object{objs[0]}
	for _, o := range objs[1:] {
		if o.Type == tt {
			co = append(co, o)
			tool.Log("cLike - Type match -", o.Type, tt)
		} else {
			conv, ok := o.TryConvert(tt)
			if !ok {
				return false, badType("value " + o.Repr() + " could not be converted for comparison")
			}
			co = append(co, conv)
		}
	}
	return cEquals(co)
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
