package lang

import "fmt"

type ObjectType uint8

const (
	objNil ObjectType = iota
	objStr
	objInt
	objBool
)

type Object struct {
	Type  ObjectType
	StrV  string
	IntV  int
	BoolV bool
}

func (o *Object) Repr() string {
	switch o.Type {
	case objNil:
		return "Nil"
	case objStr:
		return "`" + o.StrV + "`"
	case objInt:
		return fmt.Sprint(o.IntV)
	case objBool:
		return fmt.Sprint(o.BoolV)
	}
	return "?"
}
