package main

import (
	"fmt"
	"io"
	"mohazit/lang"
	"mohazit/lib"
	"os"
)

const (
	eArgs int = 1 + iota
	eFile
	eRead
	eInterpreter
	eScript
	eCleanup
)

func main() {
	lib.Load()
	if len(os.Args) < 2 {
		fmt.Println("need input file")
		exit(eArgs)
	} else {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println(err.Error())
			exit(eFile)
		}
		s, err := io.ReadAll(f)
		if err != nil {
			fmt.Println(err.Error())
			exit(eRead)
		}
		lang.Source(string(s))
		err = lang.DoAll()
		if err != nil {
			if perr, ok := err.(*lang.ParseError); ok {
				fmt.Printf("%s:%d:%d [ERROR] %s",
					os.Args[1], perr.Where.Line, perr.Where.Col, perr.Error())
			} else {
				fmt.Println(err.Error())
			}
			exit(eScript)
		}
	}
	exit(0)
}

func exit(code int) {
	if err := lib.Cleanup(); err != nil {
		fmt.Println("-- CLEANUP ERROR --")
		fmt.Println("(this usually isn't a serious problem, but should be avoided!")
		fmt.Println(err.Error())
		if code == 0 {
			os.Exit(eCleanup)
		} else {
			os.Exit(code)
		}
	}
}
