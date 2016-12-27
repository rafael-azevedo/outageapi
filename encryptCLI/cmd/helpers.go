package cmd

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

var funcMap template.FuncMap
var projectPath = ""
var inputPath = ""
var projectBase = ""

// for testing only
var testWd = ""

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(-1)
}

// Check if a file or directory exists.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// wrapper of the os package so we can test better
func getWd() (string, error) {
	if testWd == "" {
		return os.Getwd()
	}
	return testWd, nil
}

// isEmpty checks if a given path is empty.
func isEmpty(path string) (bool, error) {
	if b, _ := exists(path); !b {
		return false, fmt.Errorf("%q path does not exist", path)
	}
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fi.IsDir() {
		f, err := os.Open(path)
		// FIX: Resource leak - f.close() should be called here by defer or is missed
		// if the err != nil branch is taken.
		defer f.Close()
		if err != nil {
			return false, err
		}
		list, _ := f.Readdir(-1)
		// f.Close() - see bug fix above
		return len(list) == 0, nil
	}
	return fi.Size() == 0, nil
}

// isDir checks if a given path is a directory.
func isDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

// dirExists checks if a path exists and is a directory.
func dirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Same as WriteToDisk but checks to see if file/directory already exists.
func safeWriteToDisk(inpath string, r io.Reader) (err error) {
	dir, _ := filepath.Split(inpath)
	ospath := filepath.FromSlash(dir)

	if ospath != "" {
		err = os.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil {
			return
		}
	}

	ex, err := exists(inpath)
	if err != nil {
		return
	}
	if ex {
		return fmt.Errorf("%v already exists", inpath)
	}

	file, err := os.Create(inpath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return
}
