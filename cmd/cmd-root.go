package cmd

import (
	"fmt"
	"os"

	"ferlab/ferlease/config"
	"github.com/spf13/cobra"
)

func generateRootCmd() *cobra.Command {
	var conf *config.Config
	var confPath string

	var rootCmd = &cobra.Command{
		Use:   "ferlease",
		Short: "Manages releases of different versions of a service in fluxcd git repo using templated orchestration",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			conf, err = config.GetConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&confPath, "config", "c", "config.yml", "Path to a yaml configuration file")
	rootCmd.MarkPersistentFlagFilename("config")

	rootCmd.AddCommand(generateReleaseCmd(conf))
	rootCmd.AddCommand(generateTeardownCmd(conf))

	return rootCmd
}

func Execute() error {
	return generateRootCmd().Execute()
}
