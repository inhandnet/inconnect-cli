# inconnect CLI Installation Guide

This document is designed for AI assistants (e.g. Claude Code) to read and execute.

## Steps

### 1. Detect Platform

Determine the OS and architecture of the current machine:

| OS      | Architecture | Binary Name              |
|---------|-------------|--------------------------|
| Linux   | x86_64      | `inconnect-linux-amd64`        |
| Linux   | arm64       | `inconnect-linux-arm64`        |
| macOS   | x86_64      | `inconnect-darwin-amd64`       |
| macOS   | arm64       | `inconnect-darwin-arm64`       |
| Windows | x86_64      | `inconnect-windows-amd64.exe`  |
| Windows | arm64       | `inconnect-windows-arm64.exe`  |

Run `uname -s` and `uname -m` to detect. On macOS, if `uname -m` returns `x86_64`, check `sysctl -n sysctl.proc_translated` — if it returns `1`, the shell is running under Rosetta 2 and the native architecture is `arm64`.

### 2. Download Binary and Checksums

Download from **GitHub Releases**:

```
https://github.com/inhandnet/inconnect-cli/releases/latest/download/{BINARY_NAME}
https://github.com/inhandnet/inconnect-cli/releases/latest/download/checksums.txt
```

### 3. Verify Checksum

The `checksums.txt` file contains SHA256 checksums in the format:

```
<hash>  <filename>
```

Verify the downloaded binary:

- macOS: `shasum -a 256 <binary>`
- Linux: `sha256sum <binary>`
- Windows (PowerShell): `(Get-FileHash <binary> -Algorithm SHA256).Hash.ToLower()`

Compare the output hash with the corresponding entry in `checksums.txt`. **Do not proceed if the checksum does not match.**

### 4. Install

#### macOS / Linux

Make the binary executable and move it to the install path:

1. Try `/usr/local/bin/inconnect` — if permission denied, use `sudo` (ask the user first)
2. If the user prefers no sudo, install to `~/.local/bin/inconnect` instead (create the directory if needed, and remind the user to add `~/.local/bin` to their PATH if it's not already there)

```bash
chmod +x <binary>
mv <binary> /usr/local/bin/inconnect
```

#### Windows

Rename the binary and move it to a directory in PATH:

```powershell
# Create install directory
New-Item -ItemType Directory -Force -Path "$env:LOCALAPPDATA\inconnect"

# Move and rename
Move-Item <binary> "$env:LOCALAPPDATA\inconnect\inconnect.exe"

# Add to user PATH (persistent, takes effect in new terminal sessions)
$currentPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($currentPath -notlike "*$env:LOCALAPPDATA\inconnect*") {
    [Environment]::SetEnvironmentVariable('Path', "$currentPath;$env:LOCALAPPDATA\inconnect", 'User')
}
```

After modifying PATH, refresh the current session or open a new terminal:

```powershell
$env:Path = [Environment]::GetEnvironmentVariable('Path', 'Machine') + ';' + [Environment]::GetEnvironmentVariable('Path', 'User')
```

### 5. Verify

Run `inconnect version` to confirm the installation succeeded.

### 6. Login

```bash
inconnect auth login            # China region (cn, default)
inconnect auth login --host us  # US region
```

This opens a browser for OAuth authorization and creates a context automatically.

Available regions:

| Region | Short name | Command                       |
|--------|-----------|-------------------------------|
| China  | `cn`      | `inconnect auth login` (default)    |
| US     | `us`      | `inconnect auth login --host us`    |
| EU     | `eu`      | `inconnect auth login --host eu`    |

You can also pass a custom domain: `inconnect auth login --host ics.example.com`.

Ask the user which region they need. After login, verify with `inconnect auth status`.

## Upgrading

Once installed, the CLI can update itself in place — no need to repeat the steps above:

```bash
inconnect update          # download and install the latest release
inconnect update --check  # only check whether a newer version exists
```

It pulls from GitHub Releases and falls back to a China mirror when GitHub is unreachable. On macOS/Linux, if the install path needs elevated permissions, run `sudo inconnect update`.
