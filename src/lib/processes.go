package lib

import (
	"bytes"
	"fmt"
	"io"
	"mohazit/lang"
	"os"
	"os/exec"
	"strings"
)

func fRun(args []*lang.Object) (*lang.Object, error) {
	var cmd string
	// var annotations []string
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need command")
	}
	cmdObj := args[0]
	if cmdObj.Type != lang.ObjStr {
		return lang.NewNil(), badType.Get("command must be a string")
	}
	cmd = cmdObj.StrV + " "
	// if len(args) >= 2 {
	// 	annotObj := args[1]
	// 	if annotObj.Type != lang.ObjStr {
	// 		return lang.NewNil(), badType.Get("annotations must be a string")
	// 	}
	// 	annotations = strings.Split(annotObj.StrV, " ")
	// }

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
	if len(ctx) > 0 {
		cmdArgs = append(cmdArgs, strings.TrimSpace(ctx))
	}
	cmdProgram, err := exec.LookPath(cmdProgramName)
	if err != nil {
		return nil, err
	}
	out := &CapturedOutput{target: os.Stderr}
	out.Quiet()
	execCmd := exec.Cmd{
		Path:   cmdProgram,
		Args:   cmdArgs,
		Stdout: out,
		Stderr: out,
	}
	if err = execCmd.Run(); err != nil {
		return nil, err
	}
	data, err := out.Data()
	if err != nil {
		return nil, err
	}
	res := strings.TrimSuffix(string(data), "\n")
	return lang.NewStr(res), err
}

type CapturedOutput struct {
	target  io.Writer
	capture bytes.Buffer
	quiet   bool
}

func (o *CapturedOutput) Target(w io.Writer) {
	o.target = w
}

func (o *CapturedOutput) Quiet() {
	o.quiet = true
}

func (o *CapturedOutput) Write(p []byte) (int, error) {
	if !o.quiet {
		o.target.Write(p)
	}
	return o.capture.Write(p)
}

func (o *CapturedOutput) Data() ([]byte, error) {
	r := bytes.NewReader(o.capture.Bytes())
	return io.ReadAll(r)
}
