// scraping is a project to collect statistics about keratoconus and the prevalence of
// corneal topography by scraping some websites. It uses nimbleway.
package main

import (
	"fmt"
	"github.com/jlewi/eyecare/go/cmd/commands"
	"github.com/jlewi/p22h/backend/pkg/logging"
	"github.com/spf13/cobra"
	"os"
)

func newRootCmd() *cobra.Command {
	var level string
	var jsonLog bool
	rootCmd := &cobra.Command{
		Short: "scraping CLI",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			_, err := logging.InitLogger(level, !jsonLog)
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&level, "level", "", "info", "The logging level.")
	rootCmd.PersistentFlags().BoolVarP(&jsonLog, "json-logs", "", false, "Enable json logging.")
	return rootCmd
}

func main() {
	rootCmd := newRootCmd()
	rootCmd.AddCommand(commands.NewScrapeCmd())
	rootCmd.AddCommand(commands.NewSearchCmd())
	rootCmd.AddCommand(commands.NewDomainsCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Command failed with error: %+v", err)
		os.Exit(1)
	}
}
