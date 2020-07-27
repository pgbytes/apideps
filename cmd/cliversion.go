package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cliVersionCmd)
}

var cliVersionCmd = &cobra.Command{
	Use:   "cliversion",
	Short: "sha1 revision used to build the program and the build time",
	Run:   showVersion,
}

func showVersion(cmd *cobra.Command, args []string) {
	if Sha1ver == "" {
		Sha1ver = "local build"
	}
	if BuildTime == "" {
		BuildTime = fmt.Sprintf(time.Now().UTC().Format("2006-01-02T15:04:05"))
	}
	fmt.Printf("sha1 version: %s\n", Sha1ver)
	fmt.Printf("build time: %s UTC\n", BuildTime)
}
