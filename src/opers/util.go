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

func invArgs(txt string) error {
	return fmt.Errorf("invalid argument: %s", txt)
}

func Init() {
	lang.Funcs = lang.FuncMap{
		"say":           fSay,
		"compare":       fCompare,
		"file-create":   fFileCreate,
		"file-delete":   fFileDelete,
		"file-rename":   fFileRename,
		"file-move":     fFileRename,
		"file-copy":     fFileCopy,
		"file-write":    fFileWrite,
		"file-append":   fFileAppend,
		"dir-create":    fDirCreate,
		"dir-delete":    fDirDelete,
		"dir-rename":    fDirRename,
		"folder-create": fDirCreate,
		"folder-delete": fDirDelete,
		"folder-rename": fDirRename,
	}
	lang.Comps = lang.CompMap{
		"equals":        cEquals,
		"=":             cEquals,
		"not-equals":    cNotEquals,
		"<>":            cNotEquals,
		"like":          cLike,
		"~":             cLike,
		"chance":        cChance,
		"file-exists":   cFileExists,
		"dir-exists":    cDirExists,
		"folder-exists": cDirExists,
	}
}
