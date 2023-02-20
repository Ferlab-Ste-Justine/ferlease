package cmd

import (
	"github.com/spf13/cobra"
)

func generateTeardownCmd(cmdVars *CmdVars) *cobra.Command {
	var teardownCmd = &cobra.Command{
		Use:   "teardown",
		Short: "teardown a release in fluxcd gitops orchestration",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return teardownCmd
}