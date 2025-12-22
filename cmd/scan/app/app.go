package app

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func funcTest(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("access denied :%s (%v)\n", path, err)
		return nil
	}

	if !d.Type().IsRegular() {
		return nil
	}
	fmt.Println("scan: ", path)
	return nil
}

func Run(args []string) {
	if len(args) == 0 {
		fmt.Println("no input path")
		return
	}

	for _, arg := range args {
		_ = filepath.WalkDir(arg, funcTest)
	}
}
