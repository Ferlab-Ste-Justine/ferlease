package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"ferlab/ferlease/config"
	"ferlab/ferlease/template"

    git "github.com/Ferlab-Ste-Justine/git-sdk"
	gogit "github.com/go-git/go-git/v5"
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

func VerifyRepoSignatures(repo *gogit.Repository, signaturesPath string) error {
	keys := []string{}
	err := filepath.Walk(signaturesPath, func(fPath string, fInfo fs.FileInfo, fErr error) error {
		if fErr != nil {
			return fErr
		}

		if fInfo.IsDir() {
			return nil
		}

		key, keyErr := os.ReadFile(fPath)
		if keyErr != nil {
			return errors.New(fmt.Sprintf("Error reading accepted signature: %s", keyErr.Error()))
		}

		keys = append(keys, string(key))

		return nil
	})

	if err != nil {
		return err
	}

	return git.VerifyTopCommit(repo, keys)
}

func SetupWorkEnv(conf *config.Config) (*gogit.Repository, *template.Orchestration) {
	exists, existsErr := PathExists(conf.RepoDir)
	AbortOnErr(existsErr)

	if exists {
		err := os.RemoveAll(conf.RepoDir)
		AbortOnErr(err)
	}

	repo, _, repErr := git.SyncGitRepo(conf.RepoDir, conf.Repo, conf.Ref, conf.GitAuth.SshKey, conf.GitAuth.KnownKey)
	AbortOnErr(repErr)

	if conf.AcceptedSignatures != "" {
		verifyErr := VerifyRepoSignatures(repo, conf.AcceptedSignatures)
		AbortOnErr(verifyErr)
	}

	tmpl := template.TemplateParameters{
		Service:     conf.Service,
		Release:     conf.Release,
		Environment: conf.Environment,
	}
	orchest, orchErr := template.LoadTemplate(conf.TemplateDirectory, &tmpl)
	AbortOnErr(orchErr)

	return repo, orchest
}