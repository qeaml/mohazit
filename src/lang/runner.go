package lang

import (
	"fmt"
	"io"
	"mohazit/tool"
	"os"
	"strings"
)

const (
	LINE_END            byte = '\n'
	COMMENT_MULTI_BEGIN      = "#:"
	COMMENT_MULTI_END        = "##"
	COMMENT_SINGLE           = "#"
)

type Runner struct {
	interpreter *interpreter
	parser      *parser
	name        string
	line        int
}

func NewRunner(name string) *Runner {
	p := &parser{}
	return &Runner{
		parser: p,
		interpreter: &interpreter{
			parser:   p,
			ctx:      ctxGlobal,
			labelMap: make(map[string][]*genStmt),
		},
		name: name,
		line: 0,
	}
}

func (r *Runner) DoFile(fn string) error {
	fd, err := os.Open(fn)
	if err != nil {
		return r.errOf(err)
	}
	src, err := io.ReadAll(fd)
	if err != nil {
		return r.errOf(err)
	}
	ctx := ""
	line := ""
	isComment := false
	for _, b := range src {
		if b == LINE_END {
			r.line++
			line = strings.TrimSpace(ctx)
			ctx = ""
			tool.Log("DoFile - Line - " + line)
			if !isComment {
				if strings.HasPrefix(line, COMMENT_MULTI_BEGIN) {
					isComment = true
				} else if line == "" || strings.HasPrefix(line, COMMENT_SINGLE) {
					// bruh
				} else {
					s, err := r.parser.ParseStatement(line)
					if err != nil {
						return r.errOf(err)
					}
					err = r.interpreter.RunStatement(s)
					if err != nil {
						return r.errOf(err)
					}
				}
			} else {
				if strings.HasSuffix(line, COMMENT_MULTI_END) {
					isComment = false
				}
			}
			tool.Log("DoFile - Comment -", isComment)
		} else {
			ctx += string(rune(b))
		}
	}
	return nil
}

func (r *Runner) errOf(err error) error {
	return fmt.Errorf("%s:%d: %s", r.name, r.line, err.Error())
}
