# gopass-json use cases

## Poor man's Vault

gopass + gopass-json gives you a locally-operated secret store with structured
JSON access — without running a server, managing leases, or paying for
infrastructure. The pattern covers most of what developer tooling actually needs
from Vault:

| Vault feature | gopass equivalent |
|---|---|
| Structured secrets (KV v2) | Multi-field gopass entries, read via `gopass-json get path .` |
| Multiple secrets in one call | `gopass-json merge path1 path2 ...` → merged JSON object |
| Secret schemas / enforcement | [gopass templates](https://github.com/gopasspw/gopass/blob/master/docs/commands/templates.md) — enforce field structure at write time |
| Access from CI/CD scripts | `gopass-json get ci/deploy .token` — one line, no grep/awk |
| Baked-in secrets for deploy | chezmoi templates with `secretJSON` — resolved at apply time, never at runtime |

### Single-call secret resolution

Instead of N subprocess calls (one per secret), fetch all secrets a pipeline
needs in one round-trip:

```bash
secrets=$(gopass-json merge api/gitlab api/registry api/deploy-key)
GITLAB_TOKEN=$(echo "$secrets" | jq -r '.gitlab.secret')
REGISTRY_TOKEN=$(echo "$secrets" | jq -r '.registry.secret')
DEPLOY_KEY=$(echo "$secrets" | jq -r '."deploy-key".secret')
```

### Chezmoi integration

Configure chezmoi to use gopass-json as its secret backend (`~/.config/chezmoi/chezmoi.toml`):

```toml
[secret]
  command = "gopass-json"
  args    = ["get"]
```

Then a single cached call in any `.tmpl` file resolves all keys from one entry:

```
{{ $env := secretJSON "api/zshenv" "." -}}
export GITLAB_TOKEN={{ index $env "GITLAB_TOKEN" }}
export TAVILY_API_KEY={{ index $env "TAVILY_API_KEY" }}
```

`secretJSON` is cached per unique args set across the entire `chezmoi apply`
run — gopass-json spawns once regardless of how many templates reference the
same entry.

### What it doesn't replace

- Dynamic secret generation (Vault PKI, database credentials with TTL)
- Multi-user / team secret sharing with RBAC (use gopass teams or a real Vault)
- Audit logs with per-token access trails
