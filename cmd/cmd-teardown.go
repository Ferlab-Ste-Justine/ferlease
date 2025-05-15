package cmd

import (
	"github.com/Ferlab-Ste-Justine/ferlease/config"
	"github.com/Ferlab-Ste-Justine/ferlease/fluxcd"
	"github.com/Ferlab-Ste-Justine/ferlease/terraform"

	git "github.com/Ferlab-Ste-Justine/git-sdk"
	"github.com/spf13/cobra"
)

func Teardown(conf *config.Config) error {
	for _, confOrch := range conf.Orchestrations {
		sshCreds, sshCredsErr := git.GetSshCredentials(confOrch.GitAuth.SshKey, confOrch.GitAuth.KnownKey, confOrch.GitAuth.User)
		if sshCredsErr != nil {
			return sshCredsErr
		}

		err := git.PushChanges(func() (*git.GitRepository, error) {					
			var repo *git.GitRepository
			var commitList []string

			if confOrch.Type == "fluxcd" {
				var orchest *fluxcd.Orchestration
				repo, orchest = SetupFluxcdWorkEnv(&confOrch, conf, sshCreds)

				commitList = RemoveFluxcdOrch(orchest, conf)
			} else {
				var orchest *terraform.Orchestration
				repo, orchest = SetupTerraformWorkEnv(&confOrch, conf, sshCreds)

				commitList = RemoveTerraformOrch(orchest, conf)
			}

			var signature *git.CommitSignatureKey
			var signatureErr error
			if confOrch.CommitSignature.Key != "" {
				signature, signatureErr = git.GetSignatureKey(confOrch.CommitSignature.Key, confOrch.CommitSignature.Passphrase)
				if signatureErr != nil {
					return nil, signatureErr
				}
			}

			changes, comErr := git.CommitFiles(
				repo, 
				commitList, 
				confOrch.CommitMessage,
				git.CommitOptions{
					Name: conf.Author.Name,
					Email: conf.Author.Email,
					SignatureKey: signature,
				},
			)
			if comErr != nil {
				return nil, comErr
			}

			if !changes {
				return nil, nil
			}

			return repo, nil
		}, confOrch.Ref, sshCreds, conf.PushRetries, conf.PushRetryInterval)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateTeardownCmd(confPath *string) *cobra.Command {
	var teardownCmd = &cobra.Command{
		Use:   "teardown",
		Short: "teardown a release in fluxcd gitops orchestration",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.GetConfig(*confPath, "teardown")
			AbortOnErr(err)

			AbortOnErr(Teardown(conf))
		},
	}

	return teardownCmd
}