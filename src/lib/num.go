package lib

import (
	"math/rand"
	"mohazit/lang"
	"strconv"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().Unix()))

func pRandom(in *lang.Object) (*lang.Object, error) {
	return &lang.Object{
		Type: lang.ObjInt,
		IntV: random.Int(),
	}, nil
}

func pLimitedRandom(in *lang.Object) (*lang.Object, error) {
	if in.Type != lang.ObjInt {
		return nil, badType.Get("bound must be an integer")
	}
	return &lang.Object{
		Type: lang.ObjInt,
		IntV: random.Intn(in.IntV),
	}, nil
}

func pAtoi(in *lang.Object) (*lang.Object, error) {
	if in.Type != lang.ObjStr {
		return nil, badType.Get("input must be a string")
	}
	n, err := strconv.Atoi(in.StrV)
	if err != nil {
		return nil, err
	}
	return lang.NewInt(n), nil
}

func pStringify(in *lang.Object) (*lang.Object, error) {
	return lang.NewStr(in.String()), nil
}
