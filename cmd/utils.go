package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/Ferlab-Ste-Justine/ferlease/config"
	"github.com/Ferlab-Ste-Justine/ferlease/fluxcd"
	"github.com/Ferlab-Ste-Justine/ferlease/kustomization"
	"github.com/Ferlab-Ste-Justine/ferlease/tplcore"

    git "github.com/Ferlab-Ste-Justine/git-sdk"
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

func VerifyRepoSignatures(repo *git.GitRepository, signaturesPath string) error {
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

func SetupFluxcdWorkEnv(confOrch *config.Orchestration, conf *config.Config, sshCreds *git.SshCredentials) (*git.GitRepository, *fluxcd.Orchestration) {
	exists, existsErr := PathExists(conf.RepoDir)
	AbortOnErr(existsErr)

	if exists {
		err := os.RemoveAll(conf.RepoDir)
		AbortOnErr(err)
	}

	repo, _, repErr := git.SyncGitRepo(conf.RepoDir, confOrch.Repo, confOrch.Ref, sshCreds)
	AbortOnErr(repErr)

	if confOrch.AcceptedSignatures != "" {
		verifyErr := VerifyRepoSignatures(repo, confOrch.AcceptedSignatures)
		AbortOnErr(verifyErr)
	}

	tmpl := tplcore.TemplateParameters{
		Service:      conf.Service,
		Release:      conf.Release,
		Environment:  conf.Environment,
		CustomParams: conf.CustomParams,
	}
	orchest, orchErr := fluxcd.LoadTemplate(confOrch.TemplateDirectory, &tmpl)
	AbortOnErr(orchErr)

	return repo, orchest
}

func ApplyFluxcdOrch(orchest *fluxcd.Orchestration, conf *config.Config) []string {
	commitList := []string{}
		
	var wErr error

	fluxcdFileName := fmt.Sprintf("%s.yml", orchest.FsConventions.Naming)
	fluxcdFilePath := path.Join(conf.RepoDir, orchest.FsConventions.FluxcdDir, fluxcdFileName)
	wErr = WriteOnFile(fluxcdFilePath, orchest.FluxcdFile)
	AbortOnErr(wErr)
	commitList = append(commitList, PathRelativeToRepo(fluxcdFilePath, conf.RepoDir))

	kusPath := path.Join(conf.RepoDir, orchest.FsConventions.FluxcdDir, "kustomization.yaml")
	kus, kusErr := kustomization.GetKustomization(kusPath)
	AbortOnErr(kusErr)

	kus.AddResource(fluxcdFileName)
	rend, rendErr := kus.Render()
	AbortOnErr(rendErr)
	wErr = WriteOnFile(kusPath, rend)
	AbortOnErr(wErr)
	commitList = append(commitList, PathRelativeToRepo(kusPath, conf.RepoDir))

	for fName, fValue := range orchest.AppFiles {
		fPath := path.Join(conf.RepoDir, orchest.FsConventions.AppsDir, orchest.FsConventions.Naming, fName)
		
		mkErr := os.MkdirAll(path.Dir(fPath), 0700)
		AbortOnErr(mkErr)
		
		wErr = WriteOnFile(fPath, fValue)
		AbortOnErr(wErr)
		commitList = append(commitList, PathRelativeToRepo(fPath, conf.RepoDir))
	}

	return commitList
}

func RemoveFluxcdOrch(orchest *fluxcd.Orchestration, conf *config.Config) []string {
	commitList := []string{}

	fluxcdFileName := fmt.Sprintf("%s.yml", orchest.FsConventions.Naming)
	fluxcdFilePath := path.Join(conf.RepoDir, orchest.FsConventions.FluxcdDir, fluxcdFileName)
	rmErr := os.Remove(fluxcdFilePath)
	AbortOnErr(rmErr)
	commitList = append(commitList, PathRelativeToRepo(fluxcdFilePath, conf.RepoDir))

	kusPath := path.Join(conf.RepoDir, orchest.FsConventions.FluxcdDir, "kustomization.yaml")
	kus, kusErr := kustomization.GetKustomization(kusPath)
	AbortOnErr(kusErr)

	kus.RemoveResource(fluxcdFileName)
	rend, rendErr := kus.Render()
	AbortOnErr(rendErr)
	wErr := WriteOnFile(kusPath, rend)
	AbortOnErr(wErr)
	commitList = append(commitList, PathRelativeToRepo(kusPath, conf.RepoDir))

	for fName, _ := range orchest.AppFiles {
		fPath := path.Join(conf.RepoDir, orchest.FsConventions.AppsDir, orchest.FsConventions.Naming, fName)
		
		rmErr := os.Remove(fPath)
		AbortOnErr(rmErr)

		commitList = append(commitList, PathRelativeToRepo(fPath, conf.RepoDir))
	}

	return commitList
}