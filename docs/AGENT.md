# Agent guide: appflowy-cli

This document complements the bundled skill (`appflowy skill print`). The CLI is **unofficial** and **not affiliated with or endorsed by AppFlowy**.

## Setup

1. Install `appflowy` (release binary or `go install`).
2. User sets `APPFLOWY_URL` and `APPFLOWY_TOKEN` (Bearer JWT from AppFlowy Cloud sign-in).
3. Run `appflowy status` — exit **0** when credentials are missing; exit **3** when a token is set but `GET /api/workspace` fails auth.

## Command map

| Area | Commands |
|------|----------|
| Workspaces | `workspaces list` |
| Folders/views | `folders list WORKSPACE_ID [--depth N] [--root-view-id UUID]` (alias: `views list`) |
| Page content | `page get WORKSPACE_ID PAGE_ID [--collab-type document]` |
| Search | `search quick-notes WORKSPACE_ID QUERY` |
| Raw API | `api PATH [--method GET] [--query k=v]` |
| Meta | `status`, `version`, `skill print`, `skill install` |

Global flags: `--dry-run`, `--write`, `--confirm write`, `--json`, `--pretty`, `--jq`, `--url`, `--http-timeout`.

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | API or runtime error |
| 2 | Usage / validation error |
| 3 | Auth failure |
| 4 | Blocked write policy |

## Writes and agents

- Never pass `--write` or `--confirm write` without explicit human approval.
- Never ask for the token in chat.
- There is **no** environment variable or config file to enable writes.

## API references (public)

- OpenAPI overview: https://github.com/AppFlowy-IO/documentations/blob/main/documentation/appflowy-cloud/openapi/README.md
- Workspace folder and collab routes are documented in AppFlowy Cloud public route listings and client HTTP modules (behaviour described in this CLI only; no code copied from AppFlowy repos).
