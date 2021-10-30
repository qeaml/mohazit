package opers

import (
	"math/rand"
	"mohazit/lang"
	"mohazit/tool"
	"time"
)

func equals(objs []*lang.Object) (bool, error) {
	for _, o := range objs {
		tool.Log(o.String())
	}
	for _, a := range objs {
		for _, b := range objs {
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

func notEquals(objs []*lang.Object) (bool, error) {
	eq, err := equals(objs)
	return !eq, err
}

func chance(objs []*lang.Object) (bool, error) {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	return r.Float64() < 0.5, nil
}
