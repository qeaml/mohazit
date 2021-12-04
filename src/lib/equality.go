package lib

import (
	"errors"
	"mohazit/lang"
)

func uObjectEquals(a *lang.Object, b *lang.Object) (bool, error) {
	if a.Type != b.Type {
		return false, nil
	}
	switch a.Type {
	case lang.ObjStr:
		return a.StrV == b.StrV, nil
	case lang.ObjInt:
		return a.IntV == b.IntV, nil
	case lang.ObjBool:
		return a.BoolV == b.BoolV, nil
	case lang.ObjInv:
		return false, errors.New("invalid object")
	}
	return true, nil
}

func cEquals(args []*lang.Object) (bool, error) {
	for _, a := range args {
		for _, b := range args[1:] {
			eq, err := uObjectEquals(a, b)
			if err != nil {
				return false, err
			}
			if !eq {
				return false, nil
			}
		}
	}
	return true, nil
}

func cNotEquals(args []*lang.Object) (bool, error) {
	eq, err := cEquals(args)
	return !eq, err
}

func cLike(args []*lang.Object) (bool, error) {
	if len(args) < 1 {
		return false, moreArgs.Get("need at least 1 argument")
	}
	// type we will try to convert other objects to
	tt := args[0].Type
	// a & b after cast to target type
	var ac *lang.Object
	var bc *lang.Object
	// used for type conversions
	var ok bool

	for _, a := range args {
		for _, b := range args[1:] {
			if a.Type == tt {
				ac = a
			} else {
				ac, ok = a.TryConvert(tt)
				if !ok {
					return false, errors.New("could not convert type")
				}
			}
			if b.Type == tt {
				bc = b
			} else {
				bc, ok = b.TryConvert(tt)
				if !ok {
					return false, errors.New("could not convert type")
				}
			}
			eq, err := uObjectEquals(ac, bc)
			if err != nil {
				return false, err
			}
			if !eq {
				return false, nil
			}
		}
	}

	return true, nil
}

func cGreater(args []*lang.Object) (bool, error) {
	if len(args) < 2 {
		return false, moreArgs.Get("need at least 2 arguments to compare")
	}
	a := args[0]
	b := args[1]
	if a.Type != b.Type {
		return false, badType.Get("both arguments must be the same type")
	}
	if a.Type != lang.ObjInt {
		return false, badType.Get("arguments are not integers, cannot compare")
	}
	return a.IntV > b.IntV, nil
}

func cLesser(args []*lang.Object) (bool, error) {
	if len(args) < 2 {
		return false, moreArgs.Get("need at least 2 arguments to compare")
	}
	a := args[0]
	b := args[1]
	if a.Type != b.Type {
		return false, badType.Get("both arguments must be the same type")
	}
	if a.Type != lang.ObjInt {
		return false, badType.Get("arguments are not integers, cannot compare")
	}
	return a.IntV < b.IntV, nil
}
