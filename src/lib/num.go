package lib

import (
	"math/rand"
	"mohazit/lang"
	"strconv"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().Unix()))

func fRandom(args []*lang.Object) (*lang.Object, error) {
	return lang.NewInt(random.Int()), nil
}

func fLimitedRandom(args []*lang.Object) (*lang.Object, error) {
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need bound")
	}
	in := args[0]
	if in.Type != lang.ObjInt {
		return nil, badType.Get("bound must be an integer")
	}
	return lang.NewInt(random.Intn(int(in.IntV()))), nil
}

func fAtoi(args []*lang.Object) (*lang.Object, error) {
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need input")
	}
	in := args[0]
	if in.Type != lang.ObjStr {
		return nil, badType.Get("input must be a string")
	}
	n, err := strconv.Atoi(in.StringV())
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
	return lang.NewInt(int(in.IntV()) + 1), nil
}

func fDec(args []*lang.Object) (*lang.Object, error) {
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need input")
	}
	in := args[0]
	if in.Type != lang.ObjInt {
		return nil, badType.Get("input must be an integer")
	}
	return lang.NewInt(int(in.IntV()) - 1), nil
}

func fNeg(args []*lang.Object) (*lang.Object, error) {
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need input")
	}
	in := args[0]
	if in.Type != lang.ObjInt {
		return nil, badType.Get("input must be an integer")
	}
	return lang.NewInt(-int(in.IntV())), nil
}
