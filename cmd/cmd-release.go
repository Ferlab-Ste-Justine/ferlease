package cmd

import (
	"ferlab/ferlease/config"
	"github.com/spf13/cobra"
)

func generateReleaseCmd(conf *config.Config) *cobra.Command {
	var releaseCmd = &cobra.Command{
		Use:   "release",
		Short: "Release a new release in fluxcd gitops orchestration",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return releaseCmd
}