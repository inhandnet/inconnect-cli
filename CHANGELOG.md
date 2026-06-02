# Changelog

All notable changes to the InConnect CLI are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
