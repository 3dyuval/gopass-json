# gopass-json

Treat your gopass vault as JSON routes. Query any secret field with `jq`.

```bash
gopass-json get infra/cloud .host
gopass-json get infra/cloud '.["api-token"]'
gopass-json list
gopass-json find cloud
```

## Install

```bash
go install github.com/3dyuval/gopass-json/cmd/gopass-json@latest
```

## Usage

| Command | Output |
|---|---|
| `gopass-json get <entry>` | All fields as JSON |
| `gopass-json get <entry> <jq-filter>` | Single field value |
| `gopass-json list [pattern]` | JSON array of matching entry paths |
| `gopass-json find <query>` | JSON array of entries matching query |
