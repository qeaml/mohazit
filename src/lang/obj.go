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
	Type ObjectType
	data reflect.Value
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

func NewStr(txt string) *Object {
	return &Object{
		Type: ObjStr,
		data: reflect.ValueOf(txt),
	}
}

func NewInt(val int) *Object {
	return &Object{
		Type: ObjInt,
		data: reflect.ValueOf(val),
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
		data: reflect.ValueOf(val),
	}
}

func (o *Object) StringV() string {
	return o.data.String()
}

func (o *Object) IntV() int {
	return int(o.data.Int())
}

func (o *Object) BoolV() bool {
	return o.data.Bool()
}

func (o *Object) Clone() *Object {
	switch o.Type {
	case ObjNil:
		return NewNil()
	case ObjStr:
		return NewStr(o.String())
	case ObjInt:
		return NewInt(o.IntV())
	case ObjBool:
		return NewBool(o.BoolV())
	}
	panic("object of invalid type: " + string(o.Type))
}

func (o *Object) Repr() string {
	return fmt.Sprintf("<%s %s>", o.Type, o)
}

func (o *Object) String() string {
	switch o.Type {
	case ObjNil:
		return "nil"
	case ObjStr:
		return o.StringV()
	case ObjInt:
		return fmt.Sprint(o.IntV())
	case ObjBool:
		return fmt.Sprint(o.BoolV())
	}
	panic("object of invalid type: " + string(o.Type))
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
		data: reflect.ValueOf(o.String()),
	}, true
}

func (o *Object) convertBool() (*Object, bool) {
	v := false
	switch o.Type {
	case ObjStr:
		v = len(o.StringV()) > 0
	case ObjInt:
		v = o.IntV() > 0
	case ObjNil:
		v = false
	default:
		return nil, false
	}
	return &Object{
		Type: ObjBool,
		data: reflect.ValueOf(v),
	}, true
}

func (o *Object) convertInt() (*Object, bool) {
	v := 0
	switch o.Type {
	case ObjStr:
		parsed, err := strconv.Atoi(o.StringV())
		if err != nil {
			return nil, false
		}
		v = parsed
	case ObjBool:
		if o.BoolV() {
			v = 1
		}
	case ObjNil:
		v = 0
	default:
		return nil, false
	}
	return &Object{
		Type: ObjInt,
		data: reflect.ValueOf(v),
	}, true
}

func (a *Object) Equals(b *Object) bool {
	if a.Type != b.Type {
		return false
	}
	switch a.Type {
	case ObjNil:
		return true // nils are always equal
	case ObjInt:
		return a.IntV() == b.IntV()
	case ObjBool:
		return a.BoolV() == b.BoolV()
	case ObjStr:
		return a.StringV() == b.StringV()
	}
	panic("object of invalid type: " + string(a.Type))
}
