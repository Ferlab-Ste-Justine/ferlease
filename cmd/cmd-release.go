package cmd

import (
	"fmt"
	"os"
	"path"

	"ferlab/ferlease/config"
	"ferlab/ferlease/kustomization"

	git "github.com/Ferlab-Ste-Justine/git-sdk"
	"github.com/spf13/cobra"
	gogit "github.com/go-git/go-git/v5"
)

func generateReleaseCmd(confPath *string) *cobra.Command {
	var releaseCmd = &cobra.Command{
		Use:   "release",
		Short: "Release a new release in fluxcd gitops orchestration",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.GetConfig(*confPath, "release")
			AbortOnErr(err)
			
			err = git.PushChanges(func() (*gogit.Repository, error) {
				repo, orchest := SetupWorkEnv(conf)

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
	
				changes, comErr := git.CommitFiles(
					repo, 
					commitList, 
					conf.CommitMessage,
					git.CommitOptions{
						Name: conf.Author.Name,
						Email: conf.Author.Email,
						SignKeyPath: conf.CommitSignature.Key,
						PassphrasePath: conf.CommitSignature.Passphrase,
					},
				)
				AbortOnErr(comErr)
	
				if !changes {
					return nil, nil
				}

				return repo, nil
			}, conf.Ref, conf.PushRetries, conf.PushRetryInterval)
			AbortOnErr(err)
		},
	}

	return releaseCmd
}