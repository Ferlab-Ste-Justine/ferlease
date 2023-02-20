package cmd

import (
	"fmt"
	"os"

	"ferlab/ferlease/config"
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

func GetRepoPath(c *config.Config) string {
	return fmt.Sprintf("%s-%s", c.Service, c.Release)
}