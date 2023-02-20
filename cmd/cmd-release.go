package cmd

import (
	"fmt"
	"os"
	"path"

	"ferlab/ferlease/kustomization"
	"github.com/spf13/cobra"
)

func generateReleaseCmd(cmdVars *CmdVars) *cobra.Command {
	var releaseCmd = &cobra.Command{
		Use:   "release",
		Short: "Release a new release in fluxcd gitops orchestration",
		Run: func(cmd *cobra.Command, args []string) {
			orchest := cmdVars.Orchest
			conf := cmdVars.Conf

			var wErr error
			repoPath := GetRepoPath(conf)

			fluxcdFileName := fmt.Sprintf("%s.yml", orchest.FsConventions.Naming)
			fluxcdFilePath := path.Join(repoPath, orchest.FsConventions.FluxcdDir, fluxcdFileName)
			wErr = WriteOnFile(fluxcdFilePath, orchest.FluxcdFile)
			AbortOnErr(wErr)

			kusPath := path.Join(repoPath, orchest.FsConventions.FluxcdDir, "kustomization.yaml")
			kus, kusErr := kustomization.GetKustomization(kusPath)
			AbortOnErr(kusErr)

			kus.AddResource(fluxcdFileName)
			rend, rendErr := kus.Render()
			AbortOnErr(rendErr)
			wErr = WriteOnFile(kusPath, rend)
			AbortOnErr(wErr)

			for fName, fValue := range orchest.AppFiles {
				fPath := path.Join(repoPath, orchest.FsConventions.AppsDir, orchest.FsConventions.Naming, fName)
				
				mkErr := os.MkdirAll(path.Dir(fPath), 0700)
				AbortOnErr(mkErr)
				
				wErr = WriteOnFile(fPath, fValue)
				AbortOnErr(wErr)
			}
		},
	}

	return releaseCmd
}