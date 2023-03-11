package cmd

import (
	"fmt"
	"os"
	"path"

	"ferlab/ferlease/config"
	"ferlab/ferlease/git"
	"ferlab/ferlease/kustomization"
	"github.com/spf13/cobra"

	gogit "github.com/go-git/go-git/v5"
)

func generateTeardownCmd(confPath *string) *cobra.Command {
	var teardownCmd = &cobra.Command{
		Use:   "teardown",
		Short: "teardown a release in fluxcd gitops orchestration",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.GetConfig(*confPath)
			AbortOnErr(err)

			git.PushChanges(func() (*gogit.Repository, error) {
				commitList := []string{}
				
				repo, orchest := SetupWorkEnv(conf)
				
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

				changes, comErr := git.CommitFiles(repo, commitList, conf.CommitMessage)
				AbortOnErr(comErr)
	
				if !changes {
					return nil, nil
				}

				return repo, nil
			}, conf.Ref, conf.PushRetries, conf.PushRetryInterval)
		},
	}

	return teardownCmd
}