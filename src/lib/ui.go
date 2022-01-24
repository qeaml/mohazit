package lib

import (
	"fmt"
	"mohazit/lang"
)

func fSay(args []*lang.Object) (*lang.Object, error) {
	for _, o := range args {
		fmt.Print(o.String(), " ")
	}
	fmt.Print("\n")
	return lang.NewNil(), nil
}

func fTypeOf(args []*lang.Object) (*lang.Object, error) {
	for _, o := range args {
		fmt.Print(o.Type.String(), " ")
	}
	fmt.Print("\n")
	return lang.NewNil(), nil
}
