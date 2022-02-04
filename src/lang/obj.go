package lang

import (
	"fmt"
	"reflect"
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
	case ObjRef:
		return "Ref"
	}
	panic("invalid object type: " + string(uint8(t)))
}

type Object struct {
	Type ObjectType
	// StrV  string
	// IntV  int
	// BoolV bool
	// RefV  string
	Data reflect.Value
}

func (o *Object) Repr() string {
	return fmt.Sprintf("<%s %s>", o.Type, o)
}

func (o *Object) String() string {
	switch o.Type {
	case ObjNil:
		return "nil"
	case ObjStr:
		return o.Data.String()
	case ObjInt:
		return fmt.Sprint(o.Data.Int())
	case ObjBool:
		return fmt.Sprint(o.Data.Bool())
	}
	panic("object of invalid type: " + string(o.Type))
}

func (o *Object) Clone() *Object {
	return &Object{
		Type: o.Type,
		Data: o.Data,
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
		Data: reflect.ValueOf(o.String()),
	}, true
}

func (o *Object) convertBool() (*Object, bool) {
	v := false
	switch o.Type {
	case ObjStr:
		v = len(o.Data.String()) > 0
	case ObjInt:
		v = o.Data.Int() > 0
	case ObjNil:
		v = false
	default:
		return nil, false
	}
	return &Object{
		Type: ObjBool,
		Data: reflect.ValueOf(v),
	}, true
}

func (o *Object) convertInt() (*Object, bool) {
	v := 0
	switch o.Type {
	case ObjStr:
		parsed, err := strconv.Atoi(o.Data.String())
		if err != nil {
			return nil, false
		}
		v = parsed
	case ObjBool:
		if o.Data.Bool() {
			v = 1
		}
	case ObjNil:
		v = 0
	default:
		return nil, false
	}
	return &Object{
		Type: ObjInt,
		Data: reflect.ValueOf(v),
	}, true
}

func NewStr(txt string) *Object {
	return &Object{
		Type: ObjStr,
		Data: reflect.ValueOf(txt),
	}
}

func NewInt(val int) *Object {
	return &Object{
		Type: ObjInt,
		Data: reflect.ValueOf(val),
	}
}

func NewNil() *Object {
	return &Object{
		Type: ObjNil,
	}
}

func NewBool(val bool) *Object {
	return &Object{
		Type: ObjBool,
		Data: reflect.ValueOf(val),
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
		return a.Data.Int() == b.Data.Int()
	case ObjBool:
		return a.Data.Bool() == b.Data.Bool()
	case ObjStr, ObjRef:
		return a.Data.String() == b.Data.String()
	}
	panic("object of invalid type: " + string(a.Type))
}
