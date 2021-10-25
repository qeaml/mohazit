package tool

import (
	"fmt"
	"os"
)

var (
	Debug bool
)

func Log(a ...interface{}) {
	if Debug {
		fmt.Print("DEBUG: ")
		fmt.Println(a...)
	}
}

func init() {
	for _, f := range os.Args {
		if f == "--debug" {
			Debug = true
			break
		}
	}
}
