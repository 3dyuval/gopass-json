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
  gopass-json get cloud/infra                  # full JSON object
  gopass-json get cloud/infra .host            # single field
  gopass-json get cloud/infra -s               # secret (first line)
  gopass-json get cloud/infra '{h:.host,t:.["api-token"]}'  # projection`,
	Args:    cobra.RangeArgs(1, 2),
	Aliases: []string{"show"},
	RunE: func(cmd *cobra.Command, args []string) error {
		entry := args[0]

		secretOnly, _ := cmd.Flags().GetBool("secret")
		filter := ""
		if secretOnly {
			filter = ".secret"
		} else if len(args) > 1 {
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

func init() {
	getCmd.Flags().BoolP("secret", "s", false, "Return only the secret (first line)")
}
