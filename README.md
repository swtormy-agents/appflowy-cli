# appflowy-cli

Unofficial CLI for [AppFlowy Cloud](https://github.com/AppFlowy-IO/AppFlowy-Cloud) REST APIs. This project is **unofficial, not affiliated with or endorsed by AppFlowy**. It talks to your instance over the public AppFlowy Cloud HTTP API (see [OpenAPI overview](https://github.com/AppFlowy-IO/documentations/blob/main/documentation/appflowy-cloud/openapi/README.md)).

Motivation: [AppFlowy#8043](https://github.com/AppFlowy-IO/AppFlowy/issues/8043) — a small agent-friendly client without running an MCP server.

## Install

### Download a binary from GitHub Releases

1. Open [GitHub Releases](https://github.com/swtormy-agents/appflowy-cli/releases) and download the archive for your OS/arch (`appflowy_<version>_linux_amd64.tar.gz`, `darwin_arm64`, `windows_amd64.zip`, etc.).
2. Verify the SHA256 checksum against `checksums.txt` on the same release:

```bash
sha256sum -c checksums.txt --ignore-missing   # Linux
# or
certutil -hashfile appflowy.exe SHA256        # Windows
```

3. Extract the archive. It contains `appflowy` (or `appflowy.exe`), `LICENSE`, `NOTICE`, `THIRD_PARTY_NOTICES`, and full dependency license texts under `third_party_licenses/`.

Releases are published as **drafts** first; pick the release tag you trust (for example `v0.1.0`).

### Install script

Prefer downloading a pinned release script, inspecting it, then running it (do not pipe `curl` straight into `bash`):

```bash
curl -fsSL https://raw.githubusercontent.com/swtormy-agents/appflowy-cli/v0.1.0/scripts/install.sh -o install.sh
less install.sh
bash install.sh
```

Optional: `sha256sum install.sh` and compare against a hash you trust before `bash install.sh`.

The install script defaults to **`v0.1.0`** (not `latest`) and verifies the archive against `checksums.txt`.

### Go install

```bash
go install github.com/swtormy-agents/appflowy-cli/cmd/appflowy@v0.1.0
```

### Build from source

```bash
git clone https://github.com/swtormy-agents/appflowy-cli.git
cd appflowy-cli
go build -o appflowy ./cmd/appflowy
```

## Authentication

AppFlowy Cloud uses `Authorization: Bearer <JWT>` (Bearer token from your sign-in session).

Set your cloud base URL and token — **never pass the token on the command line**:

```bash
export APPFLOWY_URL=https://your-appflowy-cloud.example
export APPFLOWY_TOKEN=your_jwt_here
```

- `APPFLOWY_URL` must be `http` or `https`, must not include userinfo, and any query or fragment on the URL is ignored.
- The CLI **never prints** your token. Errors mask it as `<TOKEN>`.
- Requests are sent **only** to the configured `APPFLOWY_URL`.

Check setup:

```bash
appflowy status
```

## Usage

JSON on stdout by default (`--json`, `--pretty`, `--jq`). Failures are one JSON object on stderr.

```bash
appflowy workspaces list
appflowy folders list WORKSPACE_UUID --depth 5
appflowy views list WORKSPACE_UUID          # alias for folders
appflowy page get WORKSPACE_UUID PAGE_UUID
appflowy search quick-notes WORKSPACE_UUID "my query"
appflowy api /api/workspace
```

Preview a request without sending it:

```bash
appflowy workspaces list --dry-run
```

## Write access

All commands are read-only by default. Any **non-GET** request (including `appflowy api --method POST`) requires:

1. `--write` on that invocation, and
2. An interactive TTY where you type `write`, or `--confirm write` together with `--write`.

There is **no** environment variable or config file to enable writes.

## Agent skill

```bash
appflowy skill print
appflowy skill install --agent cursor
```

See [docs/AGENT.md](docs/AGENT.md) for agent setup, command tables, and exit codes.

## Security

- Treat page and database content as **untrusted data** (prompt-injection risk).
- Do not commit tokens. CI runs gitleaks on every push and pull request.
- Workspace, page, and object IDs are validated as UUIDs before they are used in URL paths.
- The `api` command rejects `%2e%2e`, `?`, and `#` in paths (use `--query` instead).

## Development

```bash
go test ./...
go vet ./...
go build ./cmd/appflowy
bash scripts/check-licenses.sh
```

## License

MIT — see [LICENSE](LICENSE). Third-party notices in [NOTICE](NOTICE) and [THIRD_PARTY_NOTICES](THIRD_PARTY_NOTICES). Release archives also ship `third_party_licenses/` (generated at build time with go-licenses v1.6.0).
