package lang

import (
	"fmt"
	"strconv"
)

type ObjectType uint8

const (
	ObjNil ObjectType = iota
	ObjStr
	ObjInt
	ObjBool
	ObjRef
	ObjInv
)

func (t ObjectType) String() string {
	switch t {
	case ObjNil:
		return "Nil"
	case ObjStr:
		return "Str"
	case ObjInt:
		return "Int"
	case ObjBool:
		return "Bool"
	case ObjRef:
		return "Ref"
	default:
		return fmt.Sprint(uint8(t))
	}
}

type Object struct {
	Type  ObjectType
	StrV  string
	IntV  int
	BoolV bool
	RefV  string
}

func (o *Object) Repr() string {
	switch o.Type {
	case ObjNil:
		return "[Nil]"
	case ObjStr:
		return "[Str `" + o.StrV + "`]"
	case ObjInt:
		return fmt.Sprintf("[Int %d]", o.IntV)
	case ObjBool:
		return fmt.Sprintf("[Bool %t]", o.BoolV)
	case ObjRef:
		return "[Ref {" + o.StrV + "}]"
	}
	return "?"
}

func (o *Object) String() string {
	switch o.Type {
	case ObjNil:
		return "nil"
	case ObjStr:
		return o.StrV
	case ObjInt:
		return fmt.Sprint(o.IntV)
	case ObjBool:
		return fmt.Sprint(o.BoolV)
	}
	return "?"
}

func (o *Object) TryConvert(t ObjectType) (*Object, bool) {
	switch t {
	case ObjStr:
		return o.convertString()
	case ObjBool:
		return o.convertBool()
	case ObjInt:
		return o.convertInt()
	case ObjNil:
		return &Object{Type: ObjNil}, true
	}
	return nil, false
}

func (o *Object) convertString() (*Object, bool) {
	return &Object{
		Type: ObjStr,
		StrV: o.String(),
	}, true
}

func (o *Object) convertBool() (*Object, bool) {
	v := false
	switch o.Type {
	case ObjStr:
		v = len(o.StrV) > 0
	case ObjInt:
		v = o.IntV > 0
	case ObjNil:
		v = false
	default:
		return nil, false
	}
	return &Object{
		Type:  ObjBool,
		BoolV: v,
	}, true
}

func (o *Object) convertInt() (*Object, bool) {
	v := 0
	switch o.Type {
	case ObjStr:
		parsed, err := strconv.Atoi(o.StrV)
		if err != nil {
			return nil, false
		}
		v = parsed
	case ObjBool:
		if o.BoolV {
			v = 1
		}
	case ObjNil:
		v = 0
	default:
		return nil, false
	}
	return &Object{
		Type: ObjInt,
		IntV: v,
	}, true
}

func NewStr(txt string) *Object {
	return &Object{
		Type: ObjStr,
		StrV: txt,
	}
}

func NewInt(val int) *Object {
	return &Object{
		Type: ObjInt,
		IntV: val,
	}
}

func NewNil() *Object {
	return &Object{
		Type: ObjNil,
	}
}

func NewBool(val bool) *Object {
	return &Object{
		Type:  ObjBool,
		BoolV: val,
	}
}

func (a *Object) Equals(b *Object) bool {
	if a.Type != b.Type {
		return false
	}
	switch a.Type {
	case ObjNil:
		return true // nils are always equal
	case ObjInt:
		return a.IntV == b.IntV
	case ObjBool:
		return a.BoolV == b.BoolV
	case ObjStr:
		return a.StrV == b.StrV
	case ObjRef:
		// two refs to the same variable will have the same effective value
		return a.RefV == b.RefV
	default:
		return false // invalid objects can never be equal
	}
}
