package main

import (
	"fmt"
	"mohazit/lang"
	"mohazit/lib"
	"os"
	"strings"
)

func main() {
	lib.Load()
	if len(os.Args) <= 1 {
		return
	} else {
		fn := os.Args[1]
		if !strings.HasSuffix(fn, ".mhzt") {
			assumed := fn + ".mhzt"
			_, err := os.Stat(assumed)
			if !os.IsNotExist(err) {
				fn = assumed
			}
		}
		r := lang.NewRunner(fn)
		if err := r.DoFile(fn); err != nil {
			fmt.Println(err.Error())
		}
	}
	if err := lib.Cleanup(); err != nil {
		fmt.Println("-- ERR -- " + err.Error())
	}
}
