package cmd

import (
	"fmt"
	"os"

	"ferlab/ferlease/config"
	"ferlab/ferlease/git"
	"ferlab/ferlease/template"
	"github.com/spf13/cobra"
	gogit "github.com/go-git/go-git/v5"
)

func generateRootCmd() *cobra.Command {
	var orchest *template.Orchestration
	var conf *config.Config
	var repo *gogit.Repository
	var confPath string

	var rootCmd = &cobra.Command{
		Use:   "ferlease",
		Short: "Manages releases of different versions of a service in fluxcd git repo using templated orchestration",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			conf, err = config.GetConfig(confPath)
			AbortOnErr(err)

			repoPath := fmt.Sprintf("%s-%s", conf.Service, conf.Release)

			exists, existsErr := PathExists(repoPath)
			AbortOnErr(existsErr)

			if exists {
				err = os.RemoveAll(repoPath)
				AbortOnErr(err)
			}

			repo, _, err = git.SyncGitRepo(repoPath, conf.Repo, conf.RepoBranch, conf.GitSshKey, conf.GitKnownKey)
			AbortOnErr(err)

			tmpl := template.TemplateParameters{
				RepoDir: repoPath,
				Service: conf.Service,
				Release: conf.Release,
			}
			orchest, err = template.LoadTemplate(conf.TemplateDirectory, &tmpl)
			AbortOnErr(err)

			fmt.Println(*orchest)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&confPath, "config", "c", "config.yml", "Path to a yaml configuration file")
	rootCmd.MarkPersistentFlagFilename("config")

	rootCmd.AddCommand(generateReleaseCmd(conf, orchest,repo ))
	rootCmd.AddCommand(generateTeardownCmd(conf, orchest, repo))

	return rootCmd
}

func Execute() error {
	return generateRootCmd().Execute()
}
