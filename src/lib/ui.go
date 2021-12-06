package lib

import (
	"fmt"
	"mohazit/lang"
)

func fSay(args []*lang.Object, i lang.InterVar) error {
	for _, o := range args {
		fmt.Print(o.String(), " ")
	}
	fmt.Print("\n")
	return nil
}

func fTypeOf(args []*lang.Object, i lang.InterVar) error {
	for _, o := range args {
		fmt.Print(lang.TypeName(o.Type), " ")
	}
	fmt.Print("\n")
	return nil
}
