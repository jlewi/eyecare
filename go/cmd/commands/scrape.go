package commands

import (
	"encoding/json"
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/jlewi/eyecare/go/api"
	"github.com/jlewi/eyecare/go/pkg/helpers"
	"github.com/jlewi/eyecare/go/pkg/scraping"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io"
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
				f, err := os.Open(input)
				defer helpers.DeferIgnoreError(f.Close)
				if err != nil {
					return errors.Wrapf(err, "Failed to read file %v", input)
				}
				d := json.NewDecoder(f)
				sites := make([]api.Site, 0, 20)
				for {
					t := &api.Site{}
					if err := d.Decode(t); err != nil {
						if err == io.EOF {
							break
						}
						return errors.Wrapf(err, "Failed to decode site")
					}
					sites = append(sites, *t)
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
