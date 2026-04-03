package main

import (
	"context"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/spf13/cobra"
)

const version = "0.1.1"

var store gopass.Store

var root = &cobra.Command{
	Use:     "gopass-json",
	Short:   "JSON-native interface to the gopass vault",
	Long:    "gopass-json treats vault entries as routes and returns structured JSON you can query with jq.",
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			return nil
		}
		gp, err := api.New(context.Background())
		if err != nil {
			return fmt.Errorf("failed to open gopass store: %w", err)
		}
		store = gp
		return nil
	},
}

func main() {
	root.AddCommand(getCmd, listCmd, findCmd)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// transport is the vault access function. Replaced in tests via fakeVault.
var transport func(request any) ([]byte, error) = storeTransport

func storeTransport(request any) ([]byte, error) {
	ctx := context.Background()

	raw, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	var msg map[string]string
	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil, err
	}

	switch msg["type"] {
	case "getData":
		return getData(ctx, msg["entry"])
	case "query":
		return queryEntries(ctx, msg["query"])
	default:
		return nil, fmt.Errorf("unknown request type: %s", msg["type"])
	}
}

func getData(ctx context.Context, entry string) ([]byte, error) {
	sec, err := store.Get(ctx, entry, "latest")
	if err != nil {
		return nil, fmt.Errorf("entry not found: %s", entry)
	}
	result := make(map[string]string)
	result["secret"] = sec.Password()
	for _, k := range sec.Keys() {
		if v, ok := sec.Get(k); ok {
			result[k] = v
		}
	}
	return json.Marshal(result)
}

func queryEntries(ctx context.Context, query string) ([]byte, error) {
	all, err := store.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list store: %w", err)
	}

	pattern := fmt.Sprintf(".*%s.*", regexp.QuoteMeta(strings.ToLower(query)))
	re := regexp.MustCompile(pattern)

	matches := []string{}
	for _, entry := range all {
		if re.MatchString(strings.ToLower(entry)) {
			matches = append(matches, entry)
		}
	}
	return json.Marshal(matches)
}

// jqPrint pipes data through jq. When a filter is given, output is raw (no quotes).
// When printing the full object, pretty-print JSON is used.
func jqPrint(data []byte, filter string) error {
	args := []string{}
	if filter == "" {
		args = append(args, ".")
	} else {
		args = append(args, "-r", filter)
	}
	cmd := exec.Command("jq", args...)
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
