package main

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <path> [jq-filter]",
	Short: "Get a vault entry as JSON",
	Long: `Fetches a vault entry. Defaults to returning the secret (first line).
Use a jq filter to extract a specific field, or '.' for the full JSON object.

Examples:
  gopass-json get cloud/infra                  # secret (first line)
  gopass-json get cloud/infra .                # full JSON object
  gopass-json get cloud/infra .host            # single field
  gopass-json get cloud/infra '{h:.host,t:.["api-token"]}'  # projection`,
	Args:    cobra.RangeArgs(1, 2),
	Aliases: []string{"show"},
	RunE: func(cmd *cobra.Command, args []string) error {
		entry := args[0]

		filter := ".secret"
		if len(args) > 1 {
			filter = args[1]
		}

		data, err := transport(map[string]string{
			"type":  "getData",
			"entry": entry,
		})
		if err != nil {
			return err
		}

		return jqPrint(data, filter)
	},
}
