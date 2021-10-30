package opers

import (
	"fmt"
	"mohazit/lang"
	"strings"
)

func say(args []*lang.Object) error {
	elem := []string{}
	for _, o := range args {
		elem = append(elem, o.String())
	}
	fmt.Println(strings.Join(elem, " "))
	return nil
}
