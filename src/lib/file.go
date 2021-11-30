package lib

import (
	"fmt"
	"mohazit/lang"
	"os"
)

func fFileCreate(args []*lang.Object, i lang.InterVar) error {
	var fileName string
	if len(args) < 1 {
		return moreArgs("need file name")
	}
	fileObj := args[0]
	if fileObj.Type != lang.ObjStr {
		return badType("file name must be a string")
	}
	fileName = fileObj.StrV

	fmt.Printf("creating file `%s`\n", fileName)

	_, err := os.Create(fileName)
	return err
}

func fFileDelete(args []*lang.Object, i lang.InterVar) error {
	var fileName string
	if len(args) < 1 {
		return moreArgs("need file name")
	}
	fileObj := args[0]
	if fileObj.Type != lang.ObjStr {
		return badType("file name must be a string")
	}
	fileName = fileObj.StrV

	fmt.Printf("deleting file `%s`\n", fileName)

	return os.Remove(fileName)
}

func fFileRename(args []*lang.Object, i lang.InterVar) error {
	var oldName string
	var newName string
	if len(args) < 2 {
		return moreArgs("need file names")
	}
	oldObj := args[0]
	if oldObj.Type != lang.ObjStr {
		return badType("file name must be a string")
	}
	oldName = oldObj.StrV
	newObj := args[0]
	if newObj.Type != lang.ObjStr {
		return badType("file name must be a string")
	}
	newName = newObj.StrV

	fmt.Printf("renaming file `%s` to `%s`\n", oldName, newName)

	return os.Rename(oldName, newName)
}

func cFileExists(args []*lang.Object) (bool, error) {
	var fileName string
	if len(args) < 1 {
		return false, moreArgs("need file name")
	}
	fileObj := args[0]
	if fileObj.Type != lang.ObjStr {
		return false, badType("file name must be a string")
	}
	fileName = fileObj.StrV

	f, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	f.Close()
	return true, nil
}
