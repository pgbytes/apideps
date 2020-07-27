package cmd

import (
	"os"

	"github.com/panshul007/apideps/config"
	"github.com/panshul007/apideps/service"
	"github.com/spf13/cobra"
)

var depFile string

func init() {
	listDepsCmd.PersistentFlags().StringVarP(&depFile, "file", "f", "apideps.yaml", "api dependencies file path")
	fetchDepsCmd.PersistentFlags().StringVarP(&depFile, "file", "f", "apideps.yaml", "api dependencies file path")
	rootCmd.AddCommand(listDepsCmd, fetchDepsCmd)
}

var listDepsCmd = &cobra.Command{
	Use:   "list",
	Short: "list the api dependencies",
	Run:   listDeps,
}

var fetchDepsCmd = &cobra.Command{
	Use:   "get",
	Short: "fetch all the api dependencies in config file",
	Run:   fetchDeps,
}

func fetchDeps(cmd *cobra.Command, args []string) {
	log.Infof("preparing to fetch the dependencies...")
	deps, err := config.LoadDeps(depFile)
	if err != nil {
		log.Errorf("error executing command: %v", err)
		os.Exit(1)
	}
	depLoader := service.NewDepLoader(log)
	err = depLoader.FetchDeps(deps)
	if err != nil {
		log.Errorf("error executing command: %v", err)
		os.Exit(1)
	}
}

func listDeps(cmd *cobra.Command, args []string) {
	log.Infof("listing api dependencies from dep file: %s", depFile)
	err := config.ListDeps(depFile)
	if err != nil {
		log.Errorf("error while listing api dependencies: %v", err)
		os.Exit(1)
	}
}
