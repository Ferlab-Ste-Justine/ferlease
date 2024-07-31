package cmd

import (
	"github.com/Ferlab-Ste-Justine/ferlease/config"

	git "github.com/Ferlab-Ste-Justine/git-sdk"
	"github.com/spf13/cobra"
)

func generateReleaseCmd(confPath *string) *cobra.Command {
	var releaseCmd = &cobra.Command{
		Use:   "release",
		Short: "Release a new release in fluxcd gitops orchestration",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.GetConfig(*confPath, "release")
			AbortOnErr(err)

			for _, confOrch := range conf.Orchestrations {
				sshCreds, sshCredsErr := git.GetSshCredentials(confOrch.GitAuth.SshKey, confOrch.GitAuth.KnownKey)
				AbortOnErr(sshCredsErr)

				err = git.PushChanges(func() (*git.GitRepository, error) {
					repo, orchest := SetupFluxcdWorkEnv(&confOrch, conf, sshCreds)

					commitList := ApplyFluxcdOrch(orchest, conf)

					var signature *git.CommitSignatureKey
					var signatureErr error
					if confOrch.CommitSignature.Key != "" {
						signature, signatureErr = git.GetSignatureKey(confOrch.CommitSignature.Key, confOrch.CommitSignature.Passphrase)
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
				}, confOrch.Ref, sshCreds, conf.PushRetries, conf.PushRetryInterval)
				AbortOnErr(err)
			}
		},
	}

	return releaseCmd
}