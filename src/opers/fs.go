package opers

import (
	"fmt"
	"io"
	"mohazit/lang"
	"os"
	"strings"
)

func sFilename(args []*lang.Object) (string, error) {
	if len(args) < 1 {
		return "", needArgs("filename missing")
	}
	fnObj := args[0]
	if fnObj.Type != lang.ObjStr {
		return "", badType("filename must be string")
	}
	return fnObj.StrV, nil
}

func fFileCreate(args []*lang.Object) error {
	fn, err := sFilename(args)
	if err != nil {
		return err
	}
	fmt.Println("Creating " + fn)
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func fFileDelete(args []*lang.Object) error {
	fn, err := sFilename(args)
	if err != nil {
		return err
	}
	fmt.Println("Deleting " + fn)
	return os.Remove(fn)
}

func fFileRename(args []*lang.Object) error {
	fn1, err := sFilename(args)
	if err != nil {
		return err
	}
	fn2, err := sFilename(args[1:])
	if err != nil {
		return err
	}
	fmt.Println("Renaming " + fn1 + " to " + fn2)
	return os.Rename(fn1, fn2)
}

func fFileCopy(args []*lang.Object) error {
	fn1, err := sFilename(args)
	if err != nil {
		return err
	}
	fn2, err := sFilename(args[1:])
	if err != nil {
		return err
	}
	f1, err := os.Open(fn1)
	if err != nil {
		return err
	}
	defer f1.Close()
	f2, err := os.Create(fn2)
	if err != nil {
		return err
	}
	defer f2.Close()
	fmt.Println("Copying from " + fn1 + " to " + fn2)
	data, err := io.ReadAll(f1)
	if err != nil {
		return err
	}
	_, err = f2.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func fFileWrite(args []*lang.Object) error {
	fn, err := sFilename(args)
	if err != nil {
		return err
	}
	elem := []string{}
	for _, o := range args[1:] {
		elem = append(elem, o.String())
	}
	data := []byte(strings.Join(elem, " "))
	fmt.Println("Writing to " + fn)
	f1, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f1.Close()
	_, err = f1.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func fFileAppend(args []*lang.Object) error {
	fn, err := sFilename(args)
	if err != nil {
		return err
	}
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Println("Appending to " + fn)
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	for _, o := range args[1:] {
		data = append(data, ' ')
		data = append(data, []byte(o.String())...)
	}
	file, err = os.Create(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(data)
	return nil
}

func fFileMove(args []*lang.Object) error {
	return nil
}

func cFileExists(args []*lang.Object) (bool, error) {
	fn, err := sFilename(args)
	if err != nil {
		return false, err
	}
	f, err := os.Open(fn)
	f.Close()
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func fDirCreate(args []*lang.Object) error {
	dirname, err := sFilename(args)
	if err != nil {
		return err
	}
	return os.MkdirAll(dirname, os.ModePerm)
}

func fDirDelete(args []*lang.Object) error {
	dirname, err := sFilename(args)
	if err != nil {
		return err
	}
	return os.Remove(dirname)
}

func fDirRename(args []*lang.Object) error {
	oldname, err := sFilename(args)
	if err != nil {
		return err
	}
	newname, err := sFilename(args[1:])
	if err != nil {
		return err
	}
	return os.Rename(oldname, newname)
}

func cDirExists(args []*lang.Object) (bool, error) {
	dirname, err := sFilename(args)
	if err != nil {
		return false, err
	}
	dirinfo, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return dirinfo.IsDir(), nil
}
