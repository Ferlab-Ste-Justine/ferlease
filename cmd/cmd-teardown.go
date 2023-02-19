package cmd

import (
	"ferlab/ferlease/config"
	"github.com/spf13/cobra"
)

func generateTeardownCmd(conf *config.Config) *cobra.Command {
	var teardownCmd = &cobra.Command{
		Use:   "teardown",
		Short: "teardown a release in fluxcd gitops orchestration",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return teardownCmd
}