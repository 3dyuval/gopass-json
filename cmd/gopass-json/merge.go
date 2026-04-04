package main

import (
	"encoding/json"
	"path"

	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge <path> [path...]",
	Short: "Fetch multiple vault entries and merge into one JSON object",
	Long: `Fetches multiple vault entries and returns them as a single JSON object
keyed by the basename of each path.

Examples:
  gopass-json merge api/gitlab.com api/tavily api/vuetify
  gopass-json merge api/gitlab.com api/tavily | jq '.tavily.secret'`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		merged, err := mergeEntries(args)
		if err != nil {
			return err
		}
		return jqPrint(merged, "")
	},
}

func mergeEntries(paths []string) ([]byte, error) {
	result := make(map[string]json.RawMessage)
	for _, entry := range paths {
		data, err := transport(map[string]string{
			"type":  "getData",
			"entry": entry,
		})
		if err != nil {
			return nil, err
		}
		result[path.Base(entry)] = json.RawMessage(data)
	}
	return json.Marshal(result)
}
