package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/jlewi/eyecare/go/pkg/helpers"
	"github.com/jlewi/eyecare/go/pkg/nimbleway"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func NewSearchCmd() *cobra.Command {
	var secretsFile string
	var query string
	var outFile string
	cmd := &cobra.Command{
		Short: "search --query=<>",
		Use:   "search",
		Run: func(cmd *cobra.Command, args []string) {
			log := zapr.NewLogger(zap.L())
			err := func() error {
				contents, err := os.ReadFile(secretsFile)
				if err != nil {
					return errors.Wrapf(err, "Failed to read file %v", secretsFile)
				}
				secrets := &nimbleway.SecretsFile{}
				if err := json.Unmarshal(contents, secrets); err != nil {
					return errors.Wrapf(err, "Failed to unmarshal file %v", secretsFile)
				}
				c := &nimbleway.Client{
					Secrets: *secrets,
				}

				if err := c.GetToken(); err != nil {
					return errors.Wrapf(err, "Failed to get nimbleway token")
				}

				log.Info("Setting query", "query", query)
				req := nimbleway.NewRequest(query)
				results, err := c.Serp(context.Background(), req)
				if err != nil {
					return errors.Wrapf(err, "SERP request failed")
				}

				sites, err := nimbleway.ToSites(results)
				if err != nil {
					return errors.Wrapf(err, "Failed to convert SERPResults to a list of sites")
				}

				// Emit it as JSONL
				f, err := os.Create(outFile)
				defer helpers.DeferIgnoreError(f.Close)
				if err != nil {
					return errors.Wrapf(err, "Failed to create file: %v", outFile)
				}
				log.Info("Writing results to file", "file", outFile)

				e := json.NewEncoder(f)
				for _, t := range sites {
					if err := e.Encode(t); err != nil {
						return errors.Wrapf(err, "Failed to encode the site")
					}
				}
				return nil
			}()
			if err != nil {
				fmt.Fprintf(os.Stdout, "Failed with error:\n%+v", err)
			}
		},
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get homedirectory; only default output directory is affected")
		homeDir = "/"
	}
	defaultSecretsFile := filepath.Join(homeDir, "secrets", "nimbleway.json")

	cmd.Flags().StringVarP(&secretsFile, "secrets", "", defaultSecretsFile, "The file containing the nimbleway username and password.")
	cmd.Flags().StringVarP(&query, "query", "q", "optometrist in san mateo", "The query to send to Google")
	cmd.Flags().StringVarP(&outFile, "outfile", "o", "results/results.jsonl", "The file to write the results to")
	return cmd
}
