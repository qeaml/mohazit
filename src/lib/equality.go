package lib

import (
	"errors"
	"mohazit/lang"
)

func cEquals(a *lang.Object, b *lang.Object) (bool, error) {
	return a.Equals(b), nil
}

func cNotEquals(a *lang.Object, b *lang.Object) (bool, error) {
	return !a.Equals(b), err
}

func cLike(a *lang.Object, b *lang.Object) (bool, error) {
	// type we will try to convert other objects to
	tt := a.Type
	// a & b after cast to target type
	var ac *lang.Object
	var bc *lang.Object
	// used for type conversions
	var ok bool

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

	return ac.Equals(bc), nil
}

func cGreater(a *lang.Object, b *lang.Object) (bool, error) {
	if a.Type != b.Type {
		return false, badType.Get("both arguments must be the same type")
	}
	if a.Type != lang.ObjInt {
		return false, badType.Get("arguments are not integers, cannot compare")
	}
	return a.IntV > b.IntV, nil
}

func cLesser(a *lang.Object, b *lang.Object) (bool, error) {
	if a.Type != b.Type {
		return false, badType.Get("both arguments must be the same type")
	}
	if a.Type != lang.ObjInt {
		return false, badType.Get("arguments are not integers, cannot compare")
	}
	return a.IntV < b.IntV, nil
}
