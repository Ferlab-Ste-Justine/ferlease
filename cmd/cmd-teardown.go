package cmd

import (
	"ferlab/ferlease/config"
	"ferlab/ferlease/template"
	"github.com/spf13/cobra"
	gogit "github.com/go-git/go-git/v5"
)

func generateTeardownCmd(conf *config.Config, orchest *template.Orchestration, repo *gogit.Repository) *cobra.Command {
	var teardownCmd = &cobra.Command{
		Use:   "teardown",
		Short: "teardown a release in fluxcd gitops orchestration",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return teardownCmd
}