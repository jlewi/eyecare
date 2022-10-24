package commands

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/jlewi/eyecare/go/pkg/scraping"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func NewScrapeCmd() *cobra.Command {
	var output string
	var input string
	var force bool
	cmd := &cobra.Command{
		Use:   "scrape",
		Short: "scrape -i <sites list>",
		Run: func(cmd *cobra.Command, args []string) {
			log := zapr.NewLogger(zap.L())
			err := func() error {
				sites, err := scraping.ReadFile(input)
				if err != nil {
					return err
				}
				log.Info("Read sites", "count", len(sites), "file", input)
				scraper := &scraping.Scraper{
					Log:       log,
					OutputDir: output,
					Force:     force,
				}

				if err := scraper.Sites(sites); err != nil {
					return err
				}
				return nil
			}()
			if err != nil {
				fmt.Fprintf(os.Stdout, "Scraping failed with error:\n%+v", err)
			}
		},
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get homedirectory; only default output directory is affected")
		homeDir = "/"
	}
	defaultDir := filepath.Join(homeDir, "scraped_sites")
	cmd.Flags().StringVarP(&input, "input", "i", "results/results.jsonl", "The JSONL file containing a list of sites to scrape. Should have been produced by the scrape command.")
	cmd.Flags().StringVarP(&output, "output", "o", defaultDir, "The directory to save the scrapings to.")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force sites to be rescraped even if they've been already scraped.")
	return cmd
}
