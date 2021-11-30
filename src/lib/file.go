package lib

import (
	"fmt"
	"io/fs"
	"mohazit/lang"
	"os"
	"sort"
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

func fWalk(args []*lang.Object, i lang.InterVar) error {
	var fileName string
	if len(args) < 1 {
		return moreArgs("need file name")
	}
	fileObj := args[0]
	if fileObj.Type != lang.ObjStr {
		return badType("file name must be a string")
	}
	fileName = fileObj.StrV

	fmt.Printf("changing working directory to `%s`\n", fileName)

	err := os.Chdir(fileName)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Printf("working directory is now `%s`\n", wd)
	return nil
}

func fFileList(args []*lang.Object, i lang.InterVar) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Printf("Directory of %s\n", wd)
	entries, err := os.ReadDir(".")
	if err != nil {
		return err
	}
	dirs := []fs.DirEntry{}
	files := []fs.DirEntry{}
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e)
		} else {
			files = append(files, e)
		}
	}
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})
	for _, d := range dirs {
		//  <directory name>/                                (DIR)      00:00:00 00.00.0000
		// ================================================================================
		niceName := " " + d.Name() + "/"
		if len(niceName) > 50 {
			niceName = niceName[:48] + "..."
		}
		for len(niceName) < 50 {
			niceName += " "
		}
		di, err := d.Info()
		if err != nil {
			return err
		}
		niceTime := di.ModTime().Format("15:04:05 02.01.2006")
		fmt.Printf("%s(DIR)      %s\n", niceName, niceTime)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	for _, f := range files {
		//  <file name>                                      <size>     00:00:00 00.00.0000
		// ================================================================================
		niceName := " " + f.Name()
		if len(niceName) > 50 {
			niceName = niceName[:48] + "..."
		}
		for len(niceName) < 50 {
			niceName += " "
		}
		fi, err := f.Info()
		if err != nil {
			return err
		}
		niceTime := fi.ModTime().Format("15:04:05 02.01.2006")
		niceSize := humanSize(fi.Size())
		for len(niceSize) < 10 {
			niceSize += " "
		}
		fmt.Printf("%s%s %s\n", niceName, niceSize, niceTime)
	}
	if len(dirs) == 1 {
		fmt.Print("1 dir, ")
	} else {
		fmt.Printf("%d dirs, ", len(dirs))
	}
	if len(files) == 1 {
		fmt.Print("1 file\n")
	} else {
		fmt.Printf("%d files\n", len(files))
	}
	return nil
}

func humanSize(size int64) string {
	var sizeF float64 = float64(size)
	var kilo float64 = 1000
	var mega float64 = kilo * 1000
	var giga float64 = mega * 1000
	var tera float64 = giga * 1000
	if sizeF >= tera {
		return fmt.Sprintf("%.2f TB", sizeF/tera)
	}
	if sizeF >= giga {
		return fmt.Sprintf("%.2f GB", sizeF/giga)
	}
	if sizeF >= mega {
		return fmt.Sprintf("%.2f MB", sizeF/mega)
	}
	if sizeF >= kilo {
		return fmt.Sprintf("%.2f kB", sizeF/kilo)
	}
	return fmt.Sprintf("%d B", size)
}