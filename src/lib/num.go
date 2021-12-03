package lib

import (
	"math/rand"
	"mohazit/lang"
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
		return nil, badType("bound must be an integer")
	}
	return &lang.Object{
		Type: lang.ObjInt,
		IntV: random.Intn(in.IntV),
	}, nil
}
