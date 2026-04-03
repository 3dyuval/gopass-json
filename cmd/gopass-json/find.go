package main

import (
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find <query>",
	Short: "Search vault entries by name",
	Long: `Search for vault entries whose path matches the query string.

Examples:
  gopass-json find api              # entries with "api" in path
  gopass-json find github           # entries with "github" in path
  gopass-json find home | jq '.[]'  # one per line`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := transport(map[string]string{
			"type":  "query",
			"query": args[0],
		})
		if err != nil {
			return err
		}

		return jqPrint(data, ".")
	},
}
