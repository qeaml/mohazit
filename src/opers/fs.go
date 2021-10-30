package opers

import (
	"fmt"
	"io"
	"mohazit/lang"
	"os"
	"strings"
)

func fn(args []*lang.Object) (string, error) {
	if len(args) < 1 {
		return "", needArgs("filename missing")
	}
	fnObj := args[0]
	if fnObj.Type != lang.ObjStr {
		return "", badType("filename must be string")
	}
	return fnObj.StrV, nil
}

func fileCreate(args []*lang.Object) error {
	fn, err := fn(args)
	if err != nil {
		return err
	}
	fmt.Println("Creating " + fn)
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func fileDelete(args []*lang.Object) error {
	fn, err := fn(args)
	if err != nil {
		return err
	}
	fmt.Println("Deleting " + fn)
	return os.Remove(fn)
}

func fileRename(args []*lang.Object) error {
	fn1, err := fn(args)
	if err != nil {
		return err
	}
	fn2, err := fn(args[1:])
	if err != nil {
		return err
	}
	fmt.Println("Renaming " + fn1 + " to " + fn2)
	return os.Rename(fn1, fn2)
}

func fileCopy(args []*lang.Object) error {
	fn1, err := fn(args)
	if err != nil {
		return err
	}
	fn2, err := fn(args[1:])
	if err != nil {
		return err
	}
	f1, err := os.Open(fn1)
	if err != nil {
		return err
	}
	f2, err := os.Create(fn2)
	if err != nil {
		return err
	}
	fmt.Println("Copying from " + fn1 + " to " + fn2)
	data, err := io.ReadAll(f1)
	if err != nil {
		return err
	}
	_, err = f2.Write(data)
	if err != nil {
		return err
	}
	f1.Close()
	f2.Close()
	return nil
}

func fileWrite(args []*lang.Object) error {
	fn, err := fn(args)
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
	_, err = f1.Write(data)
	if err != nil {
		return err
	}
	f1.Close()
	return nil
}

func fileMove(args []*lang.Object) error {
	return nil
}
