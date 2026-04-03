package main

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <path> [jq-filter]",
	Short: "Get a vault entry as JSON",
	Long: `Fetches all fields of a vault entry.
Optionally apply a jq filter to extract a specific value.

Examples:
  gopass-json get infra/cloud                  # full JSON object
  gopass-json get infra/cloud .host            # single field
  gopass-json get infra/cloud .password        # password field
  gopass-json get infra/cloud '{h:.host,t:.["api-token"]}'  # projection`,
	Args:    cobra.RangeArgs(1, 2),
	Aliases: []string{"show"},
	RunE: func(cmd *cobra.Command, args []string) error {
		entry := args[0]
		filter := ""
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
