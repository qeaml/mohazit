package main

import (
	"bufio"
	"fmt"
	"io"
	"mohazit/lang"
	"mohazit/opers"
	"mohazit/tool"
	"os"
	"strings"
)

var i = lang.NewInterpreter()

func main() {
	opers.Init()
	if len(os.Args) <= 1 {
		repl()
	} else {
		fn := os.Args[1]
		if !strings.HasSuffix(fn, ".mhzt") {
			assumed := fn + ".mhzt"
			_, err := os.Stat(assumed)
			if !os.IsNotExist(err) {
				fn = assumed
			}
		}
		err := i.RunFile(fn)
		if err != nil {
			fmt.Println("error: " + err.Error())
		}
	}
}

func repl() {
	fmt.Println("Mohazit", tool.Version, tool.Iteration)
	reader := bufio.NewReader(os.Stdin)
	var err error
	var in string
	for {
		fmt.Print(i.Context() + "> ")
		in, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("\terror: " + err.Error())
			break
		}
		if len(in) >= 2 && strings.HasPrefix(strings.TrimSpace(in), "#") {
			do := strings.ToLower(strings.TrimSpace(in[1:]))
			if strings.HasPrefix(do, "q") {
				break
			}
		}
		err = i.RunLine(in)
		if err != nil {
			fmt.Println("\terror: " + err.Error())
		}
	}
}
