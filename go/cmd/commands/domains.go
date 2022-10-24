package commands

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/jlewi/eyecare/go/api"
	"github.com/jlewi/eyecare/go/pkg/scraping"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

func NewDomainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domains",
		Short: "domains <subcommand",
	}

	cmd.AddCommand(NewAddCmd())
	return cmd
}
func NewAddCmd() *cobra.Command {
	var input string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add <domain>",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log := zapr.NewLogger(zap.L())
			err := func() error {
				sites, err := scraping.ReadFile(input)
				if err != nil {
					return err
				}
				log.Info("Read sites", "count", len(sites), "file", input)

				domain := scraping.URLToDomain(args[0])
				log.Info("Checking domain", "domain", domain)
				for _, s := range sites {
					if s.Domain == domain {
						fmt.Fprintf(os.Stdout, "Domain %v is already listed; not adding it", domain)
						return nil
					}
				}
				sites = append(sites, api.Site{
					Domain:                   domain,
					IsOptometrist:            false,
					CornealTopography:        false,
					CornealTopographyMention: "",
				})

				log.Info("Writing results to file", "file", input)
				fmt.Fprintf(os.Stdout, "Domain %v added", domain)
				return scraping.WriteFile(input, sites)
			}()
			if err != nil {
				fmt.Fprintf(os.Stdout, "Scraping failed with error:\n%+v", err)
			}
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", "results/results.jsonl", "The JSONL file containing a list of sites to scrape. Should have been produced by the scrape command.")
	return cmd
}
