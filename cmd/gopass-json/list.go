package main

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [pattern]",
	Short: "List vault entries",
	Long: `List all vault entries, optionally filtered by a pattern.

Examples:
  gopass-json list                  # all entries
  gopass-json list home             # entries matching "home"
  gopass-json list home | jq '.[]'  # one per line`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		data, err := transport(map[string]string{
			"type":  "query",
			"query": query,
		})
		if err != nil {
			return err
		}

		return jqPrint(data, ".")
	},
}
