# gopass-json — JSON-native vault CLI

## Goal

A CLI that treats the gopass vault as a collection of JSON routes — replacing fragile `grep | awk` field parsing with typed JSON queries and `jq` filters.

## Motivation

Scripts that access multi-field gopass entries typically look like:

```bash
_f() { gopass show cloud/infra | grep "^$1:" | awk '{print $2}'; }
HOST=$(_f host)
TOKEN=$(_f api-token)
```

This breaks on values containing spaces, colons, or special characters, and is fragile across gopass versions and entry formats.

`gopass-json` exposes the vault as structured JSON so callers use `jq` instead.

## Usage

```bash
gopass-json get cloud/infra                        # full JSON object of all fields
gopass-json get cloud/infra .host                  # single field via jq filter
gopass-json get cloud/infra '.["api-token"]'       # field with special chars
gopass-json list                                   # ["cloud/infra", "cloud/wifi", ...]
gopass-json list cloud                            # filtered list
gopass-json find infra                             # search by name pattern
```

## Design

- **Files are routes** — vault entry paths map directly to CLI arguments
- **jq is the query language** — any valid jq filter works as the second argument to `get`
- **No grep/awk** — all field access is typed JSON

## Request types (internal)

| Type | Purpose |
|---|---|
| `getData` | Returns all fields as a `map[string]string` for an entry |
| `query` | Search entries by name pattern |

## Structure

```
cmd/gopass-json/
├── main.go          # cobra root, gopass store transport
├── get.go           # get <entry> [jq-filter] → JSON
├── list.go          # list [pattern] → JSON array
├── find.go          # find <query> → JSON array
└── main_test.go     # transport-level tests with fake vault
docs/
└── api.md           # this file
```
