package cmd

import (
	"fmt"
	"os"
	"path"
)

func AbortOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return true, err
		}

		return false, nil
	}

	return true, nil
}

func WriteOnFile(path string, content string) error {
    f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0700)
    if err != nil {
        return err
    }
   
    _, err = f.Write([]byte(content))
    if err != nil {
        return err
    }

    err = f.Close()
	return err
}

func PathRelativeToRepo(fPath string, repo string) string {
	relative := ""
	for true {
		dir := path.Dir(fPath)
		file := path.Base(fPath)
		if file == repo {
			break
		}
		
		fPath = dir
		relative = path.Join(file, relative)
	}

	return relative
}