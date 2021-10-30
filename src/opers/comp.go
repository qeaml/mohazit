package opers

import (
	"mohazit/lang"
	"mohazit/tool"
)

func equals(objs []*lang.Object) (bool, error) {
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
	for _, o := range objs {
		tool.Log(o.String())
	}
	eq, err := equals(objs)
	return !eq, err
}
