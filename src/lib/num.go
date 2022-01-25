package lib

import (
	"math/rand"
	"mohazit/lang"
	"strconv"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().Unix()))

func fRandom(args []*lang.Object) (*lang.Object, error) {
	return &lang.Object{
		Type: lang.ObjInt,
		IntV: random.Int(),
	}, nil
}

func fLimitedRandom(args []*lang.Object) (*lang.Object, error) {
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need bound")
	}
	in := args[0]
	if in.Type != lang.ObjInt {
		return nil, badType.Get("bound must be an integer")
	}
	return &lang.Object{
		Type: lang.ObjInt,
		IntV: random.Intn(in.IntV),
	}, nil
}

func fAtoi(args []*lang.Object) (*lang.Object, error) {
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need input")
	}
	in := args[0]
	if in.Type != lang.ObjStr {
		return nil, badType.Get("input must be a string")
	}
	n, err := strconv.Atoi(in.StrV)
	if err != nil {
		return nil, err
	}
	return lang.NewInt(n), nil
}

func fStringify(args []*lang.Object) (*lang.Object, error) {
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need input")
	}
	in := args[0]
	return lang.NewStr(in.String()), nil
}

func fInc(args []*lang.Object) (*lang.Object, error) {
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need input")
	}
	in := args[0]
	if in.Type != lang.ObjInt {
		return nil, badType.Get("input must be an integer")
	}
	return lang.NewInt(in.IntV + 1), nil
}

func fDec(args []*lang.Object) (*lang.Object, error) {
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need input")
	}
	in := args[0]
	if in.Type != lang.ObjInt {
		return nil, badType.Get("input must be an integer")
	}
	return lang.NewInt(in.IntV - 1), nil
}

func fNeg(args []*lang.Object) (*lang.Object, error) {
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need input")
	}
	in := args[0]
	if in.Type != lang.ObjInt {
		return nil, badType.Get("input must be an integer")
	}
	return lang.NewInt(-in.IntV), nil
}
