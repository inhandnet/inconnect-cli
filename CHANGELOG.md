# Changelog

All notable changes to the InConnect CLI are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v0.2.1 - 2026-06-04

### Changed

- `server logs` — `--tail` and `--since` are now mutually exclusive query modes,
  matching the updated API semantics: `--tail` returns the last N lines, `--since`
  returns logs from a time offset (volume capped server-side). `--tail` is only
  sent when explicitly set; passing both flags now fails locally.

## v0.2.0 - 2026-06-02

### Added

- **Diagnostics commands** for troubleshooting router connectivity across data
  sources:
  - `connection-log list` — VPN session logs (vpn-controller): who connected /
    disconnected, bytes, duration, virtual IP.
  - `vpn-event list` — VPN auth/connection events, including auth failures with a
    reason (e.g. `invalid_cert`).
  - `router connection-events` — device MQTT online/offline events (site), with
    disconnect reasons such as `timeout` / `kicked`.
  - `server logs` — stream the org's OpenVPN server Pod logs (`--tail`/`--since`,
    raw text or line-wrapped JSON); admin only.
  - `router diagnose` — one-shot aggregation of the above (plus device
    registration events) into a single block-grouped JSON report; per-source
    failures degrade gracefully to empty arrays.

## v0.1.0 - 2026-06-02

First public release of the InConnect CLI — a command-line client for managing
the InConnect (InVPN) secure remote-access platform.

### Added

- **Browser-based OAuth login** with multi-region support (China / US / EU) and
  named contexts for switching between accounts and environments
  (`auth login`, `auth status`, `config`).
- **Full platform coverage** across command groups for VPN networks, servers,
  routers, and endpoints, plus users, organizations, roles, alerts, billing,
  firmware, data usage, audit logs, registration logs, mail, banners, tasks,
  and device running-config templates.
- **Scriptable output** — JSON by default (LLM- and pipe-friendly), with
  `--output table|yaml` and built-in `--jq` filtering.
- **Raw API escape hatch** — `inconnect api` issues authenticated requests
  against any platform endpoint.
- Cross-platform static binaries for Linux, macOS, and Windows
  (amd64 / arm64), published with SHA-256 checksums.
