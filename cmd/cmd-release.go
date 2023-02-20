package cmd

import (
	"ferlab/ferlease/config"
	"ferlab/ferlease/template"
	"github.com/spf13/cobra"
	gogit "github.com/go-git/go-git/v5"
)

func generateReleaseCmd(conf *config.Config, orchest *template.Orchestration, repo *gogit.Repository) *cobra.Command {
	var releaseCmd = &cobra.Command{
		Use:   "release",
		Short: "Release a new release in fluxcd gitops orchestration",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return releaseCmd
}