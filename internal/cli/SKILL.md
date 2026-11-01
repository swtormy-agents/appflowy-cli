---
name: appflowy
description: Unofficial AppFlowy Cloud REST CLI. Read workspaces, folder trees, page JSON, quick-note search; raw GET via api. Writes need human --write + TTY confirm. Auth APPFLOWY_URL + APPFLOWY_TOKEN only.
---

# appflowy-cli skill

Unofficial — not affiliated with or endorsed by AppFlowy.

## Auth

```bash
export APPFLOWY_URL=https://your-cloud.example
export APPFLOWY_TOKEN=your_jwt
appflowy status
```

Never print or log `APPFLOWY_TOKEN`.

## Read commands

```bash
appflowy workspaces list
appflowy folders list WORKSPACE_UUID --depth 5
appflowy page get WORKSPACE_UUID PAGE_UUID
appflowy search quick-notes WORKSPACE_UUID "query text"
appflowy api /api/workspace
```

## Writes

Non-GET (including `api --method POST`) requires `--write` and typing `write` on a TTY (or `--confirm write` with `--write`). No env override.

## Output

JSON on stdout; errors JSON on stderr. Use `--pretty` or `--jq '.data'`.

## Exit codes

0 ok, 1 runtime/API, 2 usage, 3 auth, 4 write blocked.

See `docs/AGENT.md` in the repository for full tables.
