# ics CLI

Command-line tool for the **InConnect / InVPN** platform — InHand's secure remote-connectivity manager for industrial IoT routers and gateways. It supports browser-based OAuth login, multi-environment/multi-account context switching, and full management of VPN networks, servers, routers, users, and more.

The CLI is designed to be **LLM / AI-agent friendly**: it defaults to JSON output when not attached to a terminal, surfaces backend errors as non-zero exit codes, and exposes a raw `api` passthrough for endpoints that aren't wrapped yet.

## Installation

### Build from source

```bash
# Requires Go 1.24+
make build      # Output to bin/ics
make install    # Install to $GOPATH/bin
```

On Windows you can also build directly:

```powershell
go build -o ics.exe ./cmd/ics
```

### Cross-platform build

```bash
make build-all  # linux/{amd64,arm64}, darwin/{amd64,arm64}, windows/{amd64,arm64}
```

## Quick start

### 1. Login

```bash
ics auth login                          # China region (cn, default)
ics auth login --host us                # US region
ics auth login --host dev               # Dev environment
ics auth login --host ics.example.com   # Custom domain
ics auth login --context prod           # Create/update a named context
```

Login uses the OAuth 2.0 Authorization Code flow — it opens a browser for authorization and reuses the platform's SPA OAuth client. A local callback server (default `http://localhost:18920/callback`) receives the authorization code and exchanges it for a token automatically. Supported regions: `cn`, `us`, `eu`, `dev`, `beta`, or any custom domain.

### 2. Verify

```bash
ics auth status
ics auth orgs           # Organizations you belong to
ics router list
```

> **Tip — `--oid`:** most write/query endpoints scope to an organization. The web app injects the org ID globally; the CLI does not, so pass `--oid <org-id>` explicitly when a command targets org-scoped resources (you'll get an HTTP error otherwise).

## Command reference

### Authentication (`auth`)

```bash
ics auth login                          # Browser-based OAuth login
ics auth status                         # Show current auth status
ics auth logout                         # Log out and invalidate tokens
ics auth orgs                           # List organizations you belong to
ics auth switch-org <org-id>            # Switch to another organization
ics auth register                       # Register a new organization account
ics auth impersonate --org <oid>        # Impersonate (requires ROOT privilege)
```

### Context management (`config`)

Contexts are created/updated at login via `--context`. The rest manage switching and inspection:

```bash
ics config list-contexts                # List all contexts
ics config current-context              # Show the active context
ics config use-context <name>           # Switch the active context
ics config set-context <name> ...       # Create or update a context
ics config delete-context <name>        # Delete a context
```

### Raw API calls (`api`)

For endpoints not yet wrapped, or for troubleshooting/repair:

```bash
ics api /api/invpn/networks/list                            # GET (default)
ics api "/api/devices/<id>?verbose=100"                     # GET with query
ics api /api/invpn/router/<id> -X DELETE                    # Other methods
ics api "/api/orgs/<from>/devices/<id>/transfer?to=<to>" -X PUT
```

### Routers (`router`)

The core surface for device operations. A router = a VPN-controller `Router` record bound to a site `Device`.

```bash
# Listing & details
ics router list --oid <oid>                          # List routers (default limit 20)
ics router list --online true --query IR600          # Filter by status / name+serial
ics router list --cursor 20 --limit 50               # Pagination
ics router get <id> --oid <oid>                      # Router details
ics router stats --oid <oid>                         # Online/offline counts + per-model
ics router models                                    # Supported router models
ics router locations --oid <oid>                     # Router map locations

# Lifecycle (write)
ics router create --serial <15-char-sn> --name r1 --model IR900 --oid <oid>
ics router update <id> --name new-name --subnet 10.16.42.0/24 --oid <oid>
ics router delete <id> --oid <oid>
ics router transfer <id> --to <target-oid> --oid <src-oid>   # Requires ROOT; see notes
ics router set-rip <id> --ip <real-ip> --enable --oid <oid>  # Real-IP access
ics router next-vip <id> --oid <oid>                 # Next available endpoint VIP
ics router next-subnet <id> --oid <oid>              # Next available subnet

# Remote control (online devices)
ics router exec <id> show log --oid <oid>            # Run a command, print output
ics router exec <id> "ifconfig eth0" --oid <oid>
ics router reboot <id> --oid <oid>                   # Reboot (device must ack)
ics router kick <id> --oid <oid>                     # Force-disconnect (auto-reconnects)
ics router web <id> --oid <oid>                      # Open web mgmt via ngrok tunnel
ics router web <id> --proto https --port 443 --no-browser --oid <oid>

# Configuration (three distinct sources — see table below)
ics router running-config <id> --oid <oid>           # LIVE config on the device (decoded)
ics router device-config get <id> --oid <oid>        # Platform's STORED copy
ics router device-config set <id> --content-file cfg.txt --oid <oid>  # Push a full config you supply
ics router device-config export <id> --oid <oid>     # Export stored config metadata
ics router config-send <id> --oid <oid>              # Push platform-rendered VPN-only config

# VPN client config downloads
ics router ovpn <id> --oid <oid>                     # Router OpenVPN config
ics router client-ovpn <id> --oid <oid>              # OpenVPN client config
ics router nat-conf <id> --oid <oid>                 # NAT config

# Monitoring / telemetry
ics router traffic-day <id> --oid <oid>              # Per-day traffic for a month (site)
ics router traffic-day <id> --month 2026-06 --oid <oid>
ics router online-trend <id> --oid <oid>             # Online/offline time series (default 24h)
ics router online-trend <id> --after 2026-06-01 --before 2026-06-02 --oid <oid>
ics router signal <id> --oid <oid>                   # Cellular signal: strength + quality
ics router signal <id> --fields rssi,rsrp --after 2026-06-01 --oid <oid>
```

**Which "config" command?** These four look similar but read/write different things:

| Command | Source | Direction | Notes |
|---|---|---|---|
| `running-config` | The device, **live**, decoded | read | Requires the device online. Authoritative current state. |
| `device-config get` | Platform **stored** copy | read | May differ from live. |
| `device-config set` | A config **you supply** | write | Pushes a full device config; device must be online. |
| `config-send` | Platform-**rendered** VPN config | write | Certs/CA/firewall only; can queue while offline. |

> **Transfer caveats:** `router transfer` requires a ROOT account, the target org must have a deployed VPN server, and the router's real-IP access must be **off** (`set-rip --enable=false`) or the backend rejects it (error_code 20005). The operation is **not atomic** across site + vpn-controller — read `ics router transfer --help` and exercise care.

### VPN networks (`network`)

```bash
ics network list --oid <oid>                         # List networks
ics network get <id> --oid <oid>                     # Network details
ics network create --name net1 --oid <oid>           # Create
ics network update <id> --name net2 --oid <oid>      # Update
ics network delete <id> --oid <oid>                  # Delete
ics network routers <id> --oid <oid>                 # Routers in a network
ics network accounts <id> --oid <oid>                # Accounts in a network
ics network endpoints <id> --oid <oid>               # Endpoints in a network
ics network centers <id> --oid <oid>                 # Center routers
ics network members <id> ... --oid <oid>             # Update members (routers + accounts)
```

### VPN servers (`server`)

```bash
ics server list --oid <oid>                          # List servers
ics server get <id> --oid <oid>                      # Server details
ics server create ... --oid <oid>                    # Create
ics server update <id> ... --oid <oid>               # Update
ics server delete <id> --oid <oid>                   # Delete
ics server deploy <id> --oid <oid>                   # Deploy / redeploy (K8s)
ics server stop --oid <oid>                          # Stop org servers
ics server recover --oid <oid>                       # Recover org servers
ics server command --oid <oid> ...                   # Send a command to the org's server
ics server issue-keypair --oid <oid>                 # Issue new server key pair(s)
ics server networks <id> --oid <oid>                 # Networks bound to a server
```

### VPN users (`user`)

```bash
ics user list --oid <oid>                            # List user accounts
ics user create --name u1 ... --oid <oid>            # Create
ics user update <id> ... --oid <oid>                 # Update
ics user delete <id> --oid <oid>                     # Delete
ics user lock <id> --oid <oid>                       # Lock
ics user unlock <id> --oid <oid>                     # Unlock
ics user reset-password <id> --oid <oid>             # Send password-reset email
ics user bind-mac <id> ... --oid <oid>               # Bind MAC addresses
ics user set-float-address <id> ... --oid <oid>      # Toggle floating IP
ics user issue-keypair <id> --oid <oid>              # Issue a key pair
ics user batch-issue-keypair ... --oid <oid>         # Batch issue key pairs
```

### Endpoints (`endpoint`)

```bash
ics endpoint list --oid <oid>                        # List endpoints
ics endpoint create ... --oid <oid>                  # Create on a router
ics endpoint update <id> ... --oid <oid>             # Update
ics endpoint delete <id> --oid <oid>                 # Delete
ics endpoint batch-delete ... --oid <oid>            # Batch delete
ics endpoint export --oid <oid>                      # Export to Excel
```

### Data usage (`data-usage`)

VPN data-usage statistics from vpn-controller (distinct from `router traffic-day`, which is site device-level traffic).

```bash
ics data-usage summary --month 2026-06 --oid <oid>   # Org-level summary (--month required)
ics data-usage account --oid <oid>                   # Daily per-account usage
ics data-usage account-month --oid <oid>             # Monthly per-account usage
ics data-usage router --oid <oid>                    # Daily per-router usage
ics data-usage router-month --oid <oid>              # Monthly per-router usage
ics data-usage account-export --oid <oid>            # Export account usage to Excel
ics data-usage router-export --oid <oid>             # Export router usage to Excel
```

### Alerts (`alert`)

```bash
ics alert list --oid <oid>                           # List alerts
ics alert get <id> --oid <oid>                       # Alert details
ics alert ack <id>... --oid <oid>                    # Acknowledge alerts
ics alert stats --oid <oid>                          # Acknowledgement statistics
ics alert rule list --oid <oid>                      # Manage alert rules (list/get/create/...)
```

### Config templates (`drc`)

```bash
ics drc list --oid <oid>                             # List templates
ics drc get <id> --oid <oid>                         # Template details
ics drc create --name t1 --content "..." --oid <oid> # Create
ics drc update <id> ... --oid <oid>                  # Update
ics drc delete <id> --oid <oid>                      # Delete
ics drc devices <id> ... --oid <oid>                 # Manage assigned devices
```

### Firmware (`firmware`)

```bash
ics firmware list --oid <oid>                        # List firmware packages
ics firmware get <id> --oid <oid>                    # Details
ics firmware create ... --oid <oid>                  # Create a package record
ics firmware update <id> ... --oid <oid>             # Update
ics firmware delete <id> --oid <oid>                 # Delete
ics firmware upgrade <device-id> --firmware-id <id> --oid <oid>   # Upgrade a device
ics firmware devices <id> ... --oid <oid>            # Manage devices in an upgrade job
ics firmware job-stats <id> --oid <oid>              # Upgrade job statistics
```

### Tasks (`task`)

```bash
ics task list --oid <oid>                            # List tasks
ics task cancel <id> --oid <oid>                     # Cancel a task
ics task restart <id> --oid <oid>                    # Restart a task
```

### Email notifications (`mail`)

```bash
ics mail list --oid <oid>                            # List notifications
ics mail get <id> --oid <oid>                        # Details
ics mail create ... --oid <oid>                      # Create and send
ics mail records <id> --oid <oid>                    # Recipients/records
ics mail cancel <id> --oid <oid>                     # Cancel in-progress
ics mail verify ... --oid <oid>                      # Send a test/verification email
```

### Billing (`billing`)

```bash
ics billing list-orders --oid <oid>                  # List orders/transactions
ics billing download-receipt <order-id> --oid <oid>  # Download receipt PDF
ics billing update-invoice <order-id> ... --oid <oid># Update invoice status
ics billing update-status <oid> ...                  # Update org billing status
ics billing get-seller --oid <oid>                   # Order-notification email settings
ics billing update-seller ... --oid <oid>            # Update those settings
```

### Organizations (`org`)

```bash
ics org list                                         # List organizations (ROOT sees all)
ics org get <id>                                     # Org details
ics org create ...                                   # Create (admin)
ics org delete <id>                                  # Delete
ics org settings --oid <oid>                         # Current org settings
ics org update-settings ... --oid <oid>              # Update settings
ics org export                                       # Export orgs to XLSX
```

### Other

```bash
ics banner list / current / create / revoke          # System banner messages
ics audit-log list --after 2026-05-01 --before 2026-05-10 --oid <oid>
ics audit-log export --oid <oid>                     # Export audit logs to XLS
ics register-log list <device-id> --oid <oid>        # Device registration events
ics system versions                                  # Backend service versions
ics system service <name>                            # Instances of a service
```

## Common flags & conventions

These flags are shared across commands via `internal/cmdutil`, so they behave the same everywhere they appear.

### Pagination

Available on most `list` commands:

```bash
ics router list --cursor 0  --limit 20   # First page (default)
ics router list --cursor 20 --limit 50   # Skip 20, take 50
```

| Flag | Default | Meaning |
|------|---------|---------|
| `--cursor` | `0` | Offset / cursor into the result set |
| `--limit` | `20` | Maximum number of records to return |

`--page-size` and `--per-page` are accepted as hidden aliases for `--limit`.

### Sorting

```bash
ics router list --sort createdAt,desc    # field,direction
```

`--sort` takes a `field,direction` pair (`asc` / `desc`), passed straight through to the backend.

### Field verbosity (`--verbose`)

Many backend endpoints (the common-lib `verbose` convention, 1–100, higher = more fields) let you control how many fields are returned. The CLI exposes a global `--verbose` flag (default **100** = all fields) that is **automatically applied to GET requests only**. POST/PUT/DELETE are unaffected.

```bash
ics router get <id> --oid <oid>               # verbose=100 by default (all fields)
ics router get <id> --verbose 5 --oid <oid>   # Minimal fields
ics router get <id> --verbose 0 --oid <oid>   # Omit the verbose param entirely
```

- If a request already specifies `verbose` (e.g. `ics api "/path?verbose=10"`), the CLI does not override it.
- `--verbose 0` disables injection — useful to fall back to the endpoint's own default.
- **Scope:** site endpoints (`/api/devices`, traffic, signal, …) honor it — e.g. `/api/devices/{id}` returns 4 fields at `verbose=1` vs 23 at `verbose=100`. The vpn-controller `/api/invpn/*` endpoints use fixed projections and ignore it (harmless).

### Time ranges

Time filters use the unified `--after` / `--before` flags. They accept either a plain date or a full timestamp, and are normalized before sending:

```bash
ics audit-log list --after 2026-05-01 --before 2026-05-10 --oid <oid>
ics router online-trend <id> --after 2026-06-01 --before 2026-06-02 --oid <oid>
```

| Input | Interpreted as |
|-------|----------------|
| `2026-06-01` | `00:00:00` **local** time |
| `2026-06-01T15:04:05` | local time (no zone) |
| `2026-06-01T15:04:05Z` / `...+08:00` | as given (RFC 3339) |

Values are converted to **UTC RFC 3339** for most endpoints, or to **Unix seconds** for endpoints that expect epoch timestamps (e.g. `router online-trend`). Omitting a flag falls back to that endpoint's own default window (e.g. last 7 days for `router signal`, last 24h for `online-trend`). Month-based commands use `--month` instead (`YYYY-MM` or `YYYYMM`).

> Used consistently by `audit-log`, `alert`, `billing`, and the `router` telemetry commands. Prefer `--after`/`--before` over ad-hoc start/end flags.

## Output formats

Use `-o` / `--output` to choose the format. Default is `table` in a terminal and `json` when piped (LLM/agent friendly).

| Format | TTY | Piped |
|--------|-----|-------|
| `table` | Aligned table | TSV |
| `json` | Colorized pretty JSON | Compact JSON |
| `yaml` | YAML | YAML |

```bash
ics router list -o json                               # Force JSON
ics router list --jq '.[].name'                       # Filter with a jq expression
```

> **jq + content fields:** `--jq '.content'` re-encodes the string (quotes + literal `\n`). To feed config back into `device-config set --content-file`, decode it with real `jq`/`python` `json.loads` first.

## Global flags

| Flag | Purpose |
|------|---------|
| `--context <name>` | Use a specific config context for this command |
| `--oid <org-id>` | Organization ID (required for org-scoped endpoints) |
| `-o, --output <fmt>` | `json`, `table`, or `yaml` |
| `--jq <expr>` | Filter JSON output with a jq expression |
| `--verbose <1-100>` | API field verbosity for GET requests (default 100; `0` to omit) |
| `--debug` | Print config/auth/HTTP debug info to stderr |
| `-v, --version` | Show version info |

## Environment variables

| Variable | Purpose |
|----------|---------|
| `ICS_CONTEXT` | Override the current context |
| `ICS_HOST` | Override the host in the current context |
| `ICS_TOKEN` | Override the token in the current context |
| `ICS_OID` | Default organization ID |
| `ICS_VERBOSE` | Default GET field verbosity (1-100; overrides the built-in default of 100) |
| `ICS_DEBUG` | Set to any non-empty value to enable debug output |

## Configuration file

Path: `<user-config-dir>/ics/config.yaml` (e.g. `~/.config/ics/config.yaml` on Linux, `%AppData%\ics\config.yaml` on Windows).

Stores all context information (host, token, etc.) and is managed via the `ics config` subcommands.

## Error handling

The backend's shared exception handler returns **HTTP 200 with a top-level `{"error": ..., "error_code": ...}`** body for business/internal errors (only auth errors map to 403/404). The CLI detects this centrally and surfaces it as a non-zero exit code with the message, so scripts and agents can rely on exit status.

## Development

### Prerequisites

- Go 1.24+
- [golangci-lint](https://golangci-lint.run/)

### Build & test

```bash
make build       # Build to bin/ics
make build-all   # Cross-platform build
make install     # Install to GOPATH
make test        # Run tests
make fmt         # Format code (golangci-lint fmt)
make lint        # Run golangci-lint
make clean       # Clean build artifacts
```

### Project structure

```
cmd/ics/                # CLI entry point
internal/
  api/                  # OAuth, token transport & auto-refresh, REST client, body-error detection
  browser/              # Shared browser-opening helper
  build/                # Injected Version/Commit/Date
  cmd/                  # Subcommand implementations
    auth/ config/ api/  # Authentication, contexts, raw API passthrough
    router/             # Routers (lifecycle, exec, web, config, telemetry)
    network/ server/    # VPN networks & servers
    user/ endpoint/     # VPN users & endpoints
    datausage/ billing/ # Usage statistics & billing
    alert/ mail/ banner/# Alerts, email notifications, banners
    drc/ firmware/ task/ # Config templates, firmware, tasks
    org/ system/        # Organizations, backend service info
    auditlog/ registerlog/ # Audit & registration logs
  cmdutil/              # Shared flags & helpers (time-flag parsing, etc.)
  config/               # Config file I/O, context model
  factory/              # Dependency-injection factory
  iostreams/            # Terminal output & formatters (JSON/Table/YAML/jq)
```
