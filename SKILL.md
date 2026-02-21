---
name: zoho-cli
description: Interact with Zoho REST APIs (CRM, Projects, WorkDrive, Writer, Cliq) via CLI. 118 commands for records, tasks, files, documents, and messaging. Use when the user needs to query, create, update, or manage data in any Zoho product, or when automating Zoho workflows.
license: GPL-3.0
compatibility: Requires the zoho-cli binary (Go). Needs network access to Zoho APIs. Needs ZOHO_CLIENT_ID, ZOHO_CLIENT_SECRET, and ZOHO_REFRESH_TOKEN env vars (or interactive auth).
metadata:
  author: omin8tor
  version: "1.0"
---

# zoho-cli

CLI for Zoho REST APIs. Single binary, JSON to stdout, 118 commands across CRM, Projects, WorkDrive, Writer, and Cliq.

## Install

```bash
go install github.com/omin8tor/zoho-cli/cmd/zoho@latest
```

Or download a prebuilt binary from GitHub releases.

## Authentication

Set these env vars:

```bash
export ZOHO_CLIENT_ID=1000.ABC123
export ZOHO_CLIENT_SECRET=xyz789
export ZOHO_REFRESH_TOKEN=1000.refresh_token_here
export ZOHO_DC=com
```

Tokens auto-refresh. The CLI caches access tokens at `~/.config/zoho-cli/cache/` to avoid hitting Zoho's 10-refreshes-per-10-minutes rate limit.

To get a refresh token, create a "Self Client" app at https://api-console.zoho.com/, generate a code with the needed scopes, then:

```bash
zoho auth self-client --code CODE --client-id ID --client-secret SECRET
```

Verify auth works:

```bash
zoho auth status
```

## How it works

Every command outputs JSON to stdout. Errors go to stderr. Exit codes: 0=success, 1=error, 2=auth error, 3=not found, 4=validation error.

The CLI is a thin wrapper — it passes through raw Zoho API responses without transformation. What Zoho returns is what you get.

Pipe into `jq` for filtering:

```bash
zoho crm records list Contacts --fields "Full_Name,Email" | jq '.[].Email'
```

## Quick reference by product

### CRM

```bash
zoho crm records list Contacts --fields "Full_Name,Email,Phone"
zoho crm records get Contacts RECORD_ID --fields "Full_Name,Email"
zoho crm records search Deals --criteria "(Stage:equals:Closed Won)" --fields "Deal_Name,Amount"
zoho crm records create Leads --json '{"Last_Name":"Smith","Company":"Acme"}'
zoho crm records update Leads RECORD_ID --json '{"Phone":"555-1234"}'
zoho crm records delete Leads RECORD_ID
zoho crm coql --query "SELECT Full_Name, Email FROM Contacts WHERE Email LIKE '%@acme.com' LIMIT 10"
zoho crm search-global "searchterm"
```

CRM v8 requires `--fields` on read endpoints. Without it, records come back empty.

COQL needs the `ZohoCRM.coql.READ` scope (separate from general CRM scopes).

### Projects

```bash
zoho projects core list --portal PORTAL_ID
zoho projects tasks my --portal PORTAL_ID
zoho projects tasks list --portal PORTAL_ID --project PROJECT_ID
zoho projects tasks create --portal PORTAL_ID --project PROJECT_ID --name "Task name"
zoho projects issues list --portal PORTAL_ID --project PROJECT_ID
```

Every Projects command needs `--portal`. Get portal IDs from `zoho projects core list`.

### WorkDrive

```bash
zoho drive teams me
zoho drive folders list --team TEAM_ID
zoho drive files list --folder FOLDER_ID
zoho drive files search --query "keyword" --team TEAM_ID
zoho drive download FILE_ID --output ./file.pdf
zoho drive upload ./file.xlsx --folder FOLDER_ID
zoho drive share add FILE_ID --email user@company.com --role editor
```

Navigate top-down: teams -> folders -> files.

### Writer

```bash
zoho writer details DOC_ID
zoho writer fields DOC_ID
zoho writer merge DOC_ID --json '{"name":"Alice"}' --format pdf --output ./out.pdf
zoho writer read DOC_ID
zoho writer download DOC_ID --format pdf --output ./doc.pdf
```

Writer has no "list documents" endpoint. Get doc IDs from WorkDrive.

### Cliq

```bash
zoho cliq channels list
zoho cliq channels message CHANNEL_NAME --text "message here"
zoho cliq buddies message user@company.com --text "hello"
zoho cliq messages list CHAT_ID
```

## Data centers

Set `ZOHO_DC` env var: `com` (US, default), `eu`, `in`, `com.au`, `jp`, `ca`, `sa`, `uk`, `com.cn`.

## Detailed references

- [Command reference](references/commands.md) — every command, flag, and usage pattern
- [API quirks](references/api-quirks.md) — Zoho-specific gotchas the CLI handles (or that you need to know about)
