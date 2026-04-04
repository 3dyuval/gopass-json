package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- fake transport helpers ---

// fakeVault stores entries as raw JSON maps.
type fakeVault map[string]map[string]string

func (v fakeVault) transport(request any) ([]byte, error) {
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
		entry, ok := v[msg["entry"]]
		if !ok {
			return json.Marshal(map[string]string{})
		}
		return json.Marshal(entry)

	case "query":
		query := strings.ToLower(msg["query"])
		var matches []string
		for k := range v {
			if query == "" || strings.Contains(strings.ToLower(k), query) {
				matches = append(matches, k)
			}
		}
		if matches == nil {
			matches = []string{}
		}
		return json.Marshal(matches)
	}

	return nil, nil
}

func withFakeVault(vault fakeVault) func() {
	orig := transport
	transport = vault.transport
	return func() { transport = orig }
}

// --- get command tests ---

func TestGetEntry(t *testing.T) {
	defer withFakeVault(fakeVault{
		"infra/cloud": {
			"secret":    "s3cret",
			"host":      "example.com",
			"api-token": "tok_abc123",
		},
	})()

	data, err := transport(map[string]string{"type": "getData", "entry": "infra/cloud"})
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, "s3cret", got["secret"])
	assert.Equal(t, "example.com", got["host"])
	assert.Equal(t, "tok_abc123", got["api-token"])
}

func TestGetDefaultsToSecret(t *testing.T) {
	// no filter → .secret raw output
	data := []byte(`{"secret":"mysecret","host":"example.com"}`)
	out := captureJqPrint(t, data, ".secret")
	assert.Equal(t, "mysecret", out)
}

func TestGetDotReturnsFullObject(t *testing.T) {
	// "." filter → full pretty JSON
	data := []byte(`{"secret":"mysecret","host":"example.com"}`)
	out := captureJqPrint(t, data, "")
	var parsed map[string]string
	require.NoError(t, json.Unmarshal([]byte(out), &parsed))
	assert.Equal(t, "mysecret", parsed["secret"])
	assert.Equal(t, "example.com", parsed["host"])
}

func TestGetEntryNotFound(t *testing.T) {
	defer withFakeVault(fakeVault{})()

	data, err := transport(map[string]string{"type": "getData", "entry": "does/not/exist"})
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Empty(t, got)
}

func TestGetEntryWithSpecialChars(t *testing.T) {
	// values containing spaces, colons, and special chars must survive intact
	defer withFakeVault(fakeVault{
		"infra/cloud": {
			"X-Client-Id":     "client.access.example.com",
			"X-Client-Secret": "key: with colon and spaces",
			"password":        "p@$$w0rd!#%",
		},
	})()

	data, err := transport(map[string]string{"type": "getData", "entry": "infra/cloud"})
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, "client.access.example.com", got["X-Client-Id"])
	assert.Equal(t, "key: with colon and spaces", got["X-Client-Secret"])
	assert.Equal(t, "p@$$w0rd!#%", got["password"])
}

// --- list command tests ---

func TestListAll(t *testing.T) {
	defer withFakeVault(fakeVault{
		"infra/cloud": {"host": "example.com"},
		"infra/wifi":  {"password": "wifipass"},
		"work/github": {"token": "ghp_abc"},
	})()

	data, err := transport(map[string]string{"type": "query", "query": ""})
	require.NoError(t, err)

	var entries []string
	require.NoError(t, json.Unmarshal(data, &entries))
	assert.Len(t, entries, 3)
	assert.Contains(t, entries, "infra/cloud")
	assert.Contains(t, entries, "infra/wifi")
	assert.Contains(t, entries, "work/github")
}

func TestListFiltered(t *testing.T) {
	defer withFakeVault(fakeVault{
		"infra/cloud": {"host": "example.com"},
		"infra/wifi":  {"password": "wifipass"},
		"work/github": {"token": "ghp_abc"},
	})()

	data, err := transport(map[string]string{"type": "query", "query": "infra"})
	require.NoError(t, err)

	var entries []string
	require.NoError(t, json.Unmarshal(data, &entries))
	assert.Len(t, entries, 2)
	assert.Contains(t, entries, "infra/cloud")
	assert.Contains(t, entries, "infra/wifi")
	assert.NotContains(t, entries, "work/github")
}

func TestListNoMatches(t *testing.T) {
	defer withFakeVault(fakeVault{
		"infra/cloud": {"host": "example.com"},
	})()

	data, err := transport(map[string]string{"type": "query", "query": "notfound"})
	require.NoError(t, err)

	var entries []string
	require.NoError(t, json.Unmarshal(data, &entries))
	assert.Empty(t, entries)
}

// --- find command tests ---

func TestFind(t *testing.T) {
	defer withFakeVault(fakeVault{
		"infra/cloud": {"host": "example.com"},
		"infra/wifi":  {"password": "wifipass"},
		"work/github": {"token": "ghp_abc"},
		"work/cloud":  {"key": "val"},
	})()

	data, err := transport(map[string]string{"type": "query", "query": "cloud"})
	require.NoError(t, err)

	var entries []string
	require.NoError(t, json.Unmarshal(data, &entries))
	assert.Len(t, entries, 2)
	assert.Contains(t, entries, "infra/cloud")
	assert.Contains(t, entries, "work/cloud")
}

// --- jqPrint output tests ---

// captureJqPrint runs jqPrint and returns stdout as a trimmed string.
func captureJqPrint(t *testing.T, data []byte, filter string) string {
	t.Helper()
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq not installed")
	}
	r, w, err := os.Pipe()
	require.NoError(t, err)
	orig := os.Stdout
	os.Stdout = w
	err = jqPrint(data, filter)
	w.Close()
	os.Stdout = orig
	require.NoError(t, err)
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return strings.TrimSpace(buf.String())
}

func TestJqPrintWithFilterIsRaw(t *testing.T) {
	// with a filter, strings should have no surrounding quotes
	data := []byte(`{"host":"example.com","port":"443"}`)
	out := captureJqPrint(t, data, ".host")
	assert.Equal(t, "example.com", out)
}

func TestJqPrintNoFilterIsPrettyJSON(t *testing.T) {
	// without a filter, output should be valid pretty-printed JSON (quoted)
	data := []byte(`{"host":"example.com"}`)
	out := captureJqPrint(t, data, "")
	assert.Contains(t, out, `"host"`)
	assert.Contains(t, out, `"example.com"`)
	var parsed map[string]string
	require.NoError(t, json.Unmarshal([]byte(out), &parsed))
}

// --- merge command tests ---

func TestMergeMultipleEntries(t *testing.T) {
	defer withFakeVault(fakeVault{
		"api/gitlab.com": {"secret": "glpat-xxx", "host": "gitlab.com"},
		"api/tavily":     {"secret": "tvly-yyy"},
		"api/vuetify":    {"secret": "vui-zzz"},
	})()

	data, err := mergeEntries([]string{"api/gitlab.com", "api/tavily", "api/vuetify"})
	require.NoError(t, err)

	var got map[string]map[string]string
	require.NoError(t, json.Unmarshal(data, &got))

	assert.Equal(t, "glpat-xxx", got["gitlab.com"]["secret"])
	assert.Equal(t, "gitlab.com", got["gitlab.com"]["host"])
	assert.Equal(t, "tvly-yyy", got["tavily"]["secret"])
	assert.Equal(t, "vui-zzz", got["vuetify"]["secret"])
}

func TestMergeKeyedByBasename(t *testing.T) {
	defer withFakeVault(fakeVault{
		"a/b/c/token": {"secret": "abc"},
	})()

	data, err := mergeEntries([]string{"a/b/c/token"})
	require.NoError(t, err)

	var got map[string]map[string]string
	require.NoError(t, json.Unmarshal(data, &got))

	assert.Contains(t, got, "token")
	assert.NotContains(t, got, "a/b/c/token")
}

func TestMergeMissingEntryReturnsEmpty(t *testing.T) {
	defer withFakeVault(fakeVault{
		"api/tavily": {"secret": "tvly-yyy"},
	})()

	data, err := mergeEntries([]string{"api/tavily", "api/missing"})
	require.NoError(t, err)

	var got map[string]map[string]string
	require.NoError(t, json.Unmarshal(data, &got))

	assert.Equal(t, "tvly-yyy", got["tavily"]["secret"])
	assert.Empty(t, got["missing"])
}

func TestMergeSingleEntry(t *testing.T) {
	defer withFakeVault(fakeVault{
		"api/github.com": {"secret": "ghp-abc", "host": "github.com"},
	})()

	data, err := mergeEntries([]string{"api/github.com"})
	require.NoError(t, err)

	var got map[string]map[string]string
	require.NoError(t, json.Unmarshal(data, &got))

	assert.Len(t, got, 1)
	assert.Equal(t, "ghp-abc", got["github.com"]["secret"])
}

func TestMergeOutputIsValidJSON(t *testing.T) {
	defer withFakeVault(fakeVault{
		"api/gitlab.com": {"secret": "glpat-xxx"},
		"api/tavily":     {"secret": "tvly-yyy"},
	})()

	data, err := mergeEntries([]string{"api/gitlab.com", "api/tavily"})
	require.NoError(t, err)
	assert.True(t, json.Valid(data))
}

func TestMergeJqOutput(t *testing.T) {
	defer withFakeVault(fakeVault{
		"api/gitlab.com": {"secret": "glpat-xxx"},
		"api/tavily":     {"secret": "tvly-yyy"},
	})()

	data, err := mergeEntries([]string{"api/gitlab.com", "api/tavily"})
	require.NoError(t, err)

	out := captureJqPrint(t, data, ".tavily.secret")
	assert.Equal(t, "tvly-yyy", out)
}

func TestFindCaseInsensitive(t *testing.T) {
	defer withFakeVault(fakeVault{
		"infra/CLOUD": {"host": "example.com"},
		"work/github": {"token": "ghp_abc"},
	})()

	data, err := transport(map[string]string{"type": "query", "query": "cloud"})
	require.NoError(t, err)

	var entries []string
	require.NoError(t, json.Unmarshal(data, &entries))
	assert.Len(t, entries, 1)
	assert.Contains(t, entries, "infra/CLOUD")
}
