package cmd

import (
	"github.com/panshul007/apideps/logger"

	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Sha1ver revision used to build the program
	Sha1ver string
	// BuildTime when the executable was built
	BuildTime string
	verbose   bool
	log       logger.GenericLogger
)

// RequestTimeoutSecondsDefault Duh
const RequestTimeoutSecondsDefault = 30

func init() {
	cobra.OnInitialize(setup)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable for verbose operation")
}

func setup() {
	if err := logger.SetupLogger(verbose); err != nil {
		fmt.Printf("error while setting up logger: %v", err)
	}
	log = logger.Logger()
}

var rootCmd = &cobra.Command{
	Use:   "apideps",
	Short: "API dependency manager",
}

// Execute the root command.
func Execute(sha1ver, buildTime string) {
	Sha1ver = sha1ver
	BuildTime = buildTime

	if err := rootCmd.Execute(); err != nil {
		// Exit as success if called with no arguments (same behaviour as
		// docker and other cobra based cli)
		if len(os.Args[1:]) == 0 {
			os.Exit(0)
		}
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		log.Errorf("Error while executing command: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
