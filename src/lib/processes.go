package lib

import (
	"fmt"
	"mohazit/lang"
	"os"
	"os/exec"
	"strings"
)

func fRun(args []*lang.Object, i lang.InterVar) error {
	var cmd string
	var annotations []string
	if len(args) < 1 {
		return moreArgs.Get("need command")
	}
	cmdObj := args[0]
	if cmdObj.Type != lang.ObjStr {
		return badType.Get("command must be a string")
	}
	cmd = cmdObj.StrV + " "
	if len(args) >= 2 {
		annotObj := args[1]
		if annotObj.Type != lang.ObjStr {
			return badType.Get("annotations must be a string")
		}
		annotations = strings.Split(annotObj.StrV, " ")
	}

	fmt.Printf("command: %s\n", cmd)

	cmdProgramName := ""
	cmdArgs := []string{}
	hasProg := false
	ctx := ""
	for _, c := range cmd {
		if c == ' ' {
			if !hasProg {
				cmdProgramName = strings.TrimSpace(ctx)
				cmdArgs = append(cmdArgs, cmdProgramName)
				ctx = ""
				hasProg = true
			} else {
				cmdArgs = append(cmdArgs, strings.TrimSpace(ctx))
				ctx = ""
			}
			continue
		}
		ctx += string(c)
	}
	cmdProgram, err := exec.LookPath(cmdProgramName)
	if err != nil {
		return err
	}
	uo := os.Stderr
	for _, a := range annotations {
		if strings.ToLower(a) == "quiet" {
			f, err := os.Open(os.DevNull)
			if err != nil {
				return err
			}
			uo = f
		}
	}
	execCmd := exec.Cmd{
		Path:   cmdProgram,
		Args:   cmdArgs,
		Stdout: uo,
		Stderr: uo,
	}
	return execCmd.Run()
}
