package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/Ferlab-Ste-Justine/ferlease/config"
	"github.com/Ferlab-Ste-Justine/ferlease/kustomization"

	git "github.com/Ferlab-Ste-Justine/git-sdk"
	"github.com/spf13/cobra"
)

func generateTeardownCmd(confPath *string) *cobra.Command {
	var teardownCmd = &cobra.Command{
		Use:   "teardown",
		Short: "teardown a release in fluxcd gitops orchestration",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.GetConfig(*confPath, "teardown")
			AbortOnErr(err)

			sshCreds, sshCredsErr := git.GetSshCredentials(conf.GitAuth.SshKey, conf.GitAuth.KnownKey)
			AbortOnErr(sshCredsErr)

			err = git.PushChanges(func() (*git.GitRepository, error) {
				commitList := []string{}
				
				repo, orchest := SetupWorkEnv(conf, sshCreds)
				
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

				var signature *git.CommitSignatureKey
				var signatureErr error
				if conf.CommitSignature.Key != "" {
					signature, signatureErr = git.GetSignatureKey(conf.CommitSignature.Key, conf.CommitSignature.Passphrase)
					AbortOnErr(signatureErr)
				}

				changes, comErr := git.CommitFiles(
					repo, 
					commitList, 
					conf.CommitMessage,
					git.CommitOptions{
						Name: conf.Author.Name,
						Email: conf.Author.Email,
						SignatureKey: signature,
					},
				)
				AbortOnErr(comErr)
	
				if !changes {
					return nil, nil
				}

				return repo, nil
			}, conf.Ref, sshCreds, conf.PushRetries, conf.PushRetryInterval)
			AbortOnErr(err)
		},
	}

	return teardownCmd
}