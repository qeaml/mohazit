package main

import (
	"bufio"
	"fmt"
	"mohazit/lang"
	"mohazit/lib"
	"os"
	"strings"
)

func main() {
	lib.Load()
	if len(os.Args) <= 1 {
		r := lang.NewRunner("REPL")
		input := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			line, err := input.ReadString('\n')
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(3)
			}
			if strings.HasPrefix(strings.TrimSpace(line), "#q") {
				break
			}
			if err = r.RunLine(line); err != nil {
				fmt.Println(err.Error())
			}
		}
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
			os.Exit(1)
		}
	}
	if err := lib.Cleanup(); err != nil {
		fmt.Println("-- ERR -- " + err.Error())
		os.Exit(2)
	}
}
