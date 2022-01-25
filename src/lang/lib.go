package lang

// snip!

type VFunc func([]*Object) (*Object, error)
type VComp func(*Object, *Object) (bool, error)
type VFuncMap map[string]VFunc
type VCompMap map[string]VComp

var Funcs = make(VFuncMap)
var Comps = make(VCompMap)
