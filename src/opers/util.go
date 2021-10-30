package opers

import (
	"fmt"
	"mohazit/lang"
)

func needArgs(txt string) error {
	return fmt.Errorf("not enough arguments: %s", txt)
}

func badType(txt string) error {
	return fmt.Errorf("wrong argument type: %s", txt)
}

func Init() {
	lang.Funcs = lang.FuncMap{
		"say":         say,
		"file-create": fileCreate,
		"file-delete": fileDelete,
		"file-rename": fileRename,
		"file-move":   fileMove,
		"file-copy":   fileCopy,
		"file-write":  fileWrite,
	}
	lang.Comps = lang.CompMap{
		"equals":     equals,
		"=":          equals,
		"not-equals": notEquals,
		"chance":     chance,
	}
}
