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
	}
	panic("invalid object type: " + string(uint8(t)))
}

type Object struct {
	Type  ObjectType
	StrV  string
	IntV  int
	BoolV bool
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
	}
	panic("object of invalid type: " + string(uint8(o.Type)))
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
	panic("object of invalid type: " + string(o.Type))
}

func (o *Object) Clone() *Object {
	return &Object{
		Type:  o.Type,
		StrV:  o.StrV,
		IntV:  o.IntV,
		BoolV: o.BoolV,
	}
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
	panic("object of invalid type: " + string(o.Type))
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

func NewObject(val interface{}) *Object {
	if val == nil {
		return NewNil()
	} else if v, ok := val.(*Object); ok {
		return v.Clone()
	} else if v, ok := val.(string); ok {
		return NewStr(v)
	} else if v, ok := val.(int); ok {
		return NewInt(v)
	} else if v, ok := val.(bool); ok {
		return NewBool(v)
	}
	panic("unsupported value: " + fmt.Sprint(val))
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
	}
	panic("object of invalid type: " + string(a.Type))
}
