package lib

import "fmt"

type genericError struct {
	msg   string
	code  string
	fatal bool
}

func (e *genericError) Error() string {
	return e.msg
}

func (e *genericError) Code() string {
	return e.code
}

func (e *genericError) Fatal() bool {
	return e.fatal
}

type lazyError struct {
	format string
	code   string
}

func LazyError(format, code string) *lazyError {
	return &lazyError{format, code}
}

func (g *lazyError) Get(txt string) error {
	return &genericError{fmt.Sprintf(g.format, txt), g.code, false}
}

func (g *lazyError) Fail(txt string) error {
	return &genericError{fmt.Sprintf(g.format, txt), g.code, true}
}

func (g *lazyError) Of(err error) error {
	return &genericError{fmt.Sprintf(g.format, err.Error()), g.code, true}
}

var (
	moreArgs = LazyError("function: not enough arguments: %s", "fnc_moreargs")
	badType  = LazyError("function: wrong type: %s", "fnc_badtype")
	badState = LazyError("function: unexpected: %s", "fnc_badstate")
)
