package cmd

import (
	"github.com/spf13/cobra"
)

func generateRootCmd() *cobra.Command {
	var confPath string

	var rootCmd = &cobra.Command{
		Use:   "ferlease",
		Short: "Manages releases of different versions of a service in fluxcd git repo using templated orchestration",
	}

	rootCmd.PersistentFlags().StringVarP(&confPath, "config", "c", "config.yml", "Path to a yaml configuration file")
	rootCmd.MarkPersistentFlagFilename("config")

	rootCmd.AddCommand(generateReleaseCmd(&confPath))
	rootCmd.AddCommand(generateTeardownCmd(&confPath))

	return rootCmd
}

func Execute() error {
	return generateRootCmd().Execute()
}