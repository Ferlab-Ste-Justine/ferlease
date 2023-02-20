package cmd

import (
	"os"

	"ferlab/ferlease/config"
	"ferlab/ferlease/git"
	"ferlab/ferlease/template"
	"github.com/spf13/cobra"
	gogit "github.com/go-git/go-git/v5"
)

type CmdVars struct {
	Orchest *template.Orchestration
	Conf    *config.Config
	Repo    *gogit.Repository
}

func generateRootCmd() *cobra.Command {
	var cmdVars CmdVars
	var confPath string

	var rootCmd = &cobra.Command{
		Use:   "ferlease",
		Short: "Manages releases of different versions of a service in fluxcd git repo using templated orchestration",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			conf, err := config.GetConfig(confPath)
			AbortOnErr(err)

			repoPath := GetRepoPath(conf)

			exists, existsErr := PathExists(repoPath)
			AbortOnErr(existsErr)

			if exists {
				err = os.RemoveAll(repoPath)
				AbortOnErr(err)
			}

			repo, _, repErr := git.SyncGitRepo(repoPath, conf.Repo, conf.RepoBranch, conf.GitSshKey, conf.GitKnownKey)
			AbortOnErr(repErr)

			tmpl := template.TemplateParameters{
				RepoDir: repoPath,
				Service: conf.Service,
				Release: conf.Release,
			}
			orchest, orchErr := template.LoadTemplate(conf.TemplateDirectory, &tmpl)
			AbortOnErr(orchErr)

			cmdVars.Orchest = orchest
			cmdVars.Conf = conf
			cmdVars.Repo = repo
		},
	}

	rootCmd.PersistentFlags().StringVarP(&confPath, "config", "c", "config.yml", "Path to a yaml configuration file")
	rootCmd.MarkPersistentFlagFilename("config")

	rootCmd.AddCommand(generateReleaseCmd(&cmdVars))
	rootCmd.AddCommand(generateTeardownCmd(&cmdVars))

	return rootCmd
}

func Execute() error {
	return generateRootCmd().Execute()
}