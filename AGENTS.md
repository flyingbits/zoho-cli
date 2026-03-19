# Agent Instructions

## Project: zoho-cli

CLI for Zoho REST APIs (CRM, Projects, WorkDrive, Writer, Cliq, Expense, Sheet, Mail, Books).

## Tech Stack

- Go 1.25+, urfave/cli v3 (CLI framework), stdlib net/http (HTTP)
- Task runner: mise
- Issue tracker: tk

## Commands

```bash
mise run build         # Build Go binary
mise run test          # Run Go tests
mise run test:all      # Run all Go tests including integration
mise run lint          # Run Go linter (go vet)
mise run fmt           # Format Go code (gofmt)
mise run typecheck     # Go type check (go build)
./zoho --help          # Run CLI
```

## Quality Gates (run before commit)

```bash
mise run lint && mise run typecheck && mise run test
```

## Issue Tracking (tk)

```bash
tk ready               # Find available work
tk show <id>           # View issue details
tk start <id>          # Claim work (set in_progress)
tk close <id>          # Complete work
tk ls                  # List all open issues
tk blocked             # Show blocked issues
```

## Architecture

- `cmd/zoho/main.go` - Entry point
- `internal/auth/` - OAuth flows, token management, config resolution
- `internal/http/` - HTTP client with auto-refresh, DC maps
- `internal/output/` - JSON output, --help-all schema display
- `internal/crm/` - CRM subcommands (29 commands)
- `internal/projects/` - Projects subcommands (39 commands)
- `internal/drive/` - WorkDrive subcommands (26 commands)
- `internal/writer/` - Writer subcommands (8 commands)
- `internal/cliq/` - Cliq subcommands (12 commands)
- `internal/mail/` - Mail subcommands (Organization, Domains, Groups, Users, Policy, Accounts, Folders, Labels, Messages, Signatures, Threads, Tasks, Bookmarks, Notes, Logs)
- `internal/inventory/` - Inventory subcommands (Organizations, Contacts, Contact Persons, Item Groups, Items, Composite Items, Item Adjustments, Transfer Orders, Sales Orders, Packages, Shipment Orders, Invoices, Retainer Invoices, Customer Payments, Sales Returns, Credit Notes, Purchase Orders, Purchase Receives, Bills, Vendor Credits, Locations, Price Lists, Users, Taxes, Currency, Reporting Tags)
- `internal/invoice/` - Invoice subcommands (Organizations, Items, Price Lists, Contacts, Contact Persons, Estimates, Invoices, Recurring Invoices, Customer Payments, Retainer Invoices, Credit Notes, Expenses, Recurring Expenses, Projects, Tasks, Time Entries, Users, Taxes, Expense Categories, Currency, CRM Integration)
- `internal/sheet/` - Sheet subcommands (97 commands across Workbooks, Worksheets, Tables, Records, Cells, Content, Format, Named Ranges, Merge, Premium, Utility)
- `internal/books/` - Books subcommands (Zoho Books API v3)

### Reference implementations
- `~/Projects/work/rhi/ai_agent/rhi-agent/src/zoho/` (original endpoints)

## Environment Variables

- `ZOHO_CLIENT_ID`, `ZOHO_CLIENT_SECRET`, `ZOHO_REFRESH_TOKEN`, `ZOHO_DC` - Auth (handled in internal/auth/config.go)
- `ZOHO_PORTAL_ID` - Default for `--portal` flag (Projects commands)
- `ZOHO_TEAM_ID` - Default for `--team` flag (WorkDrive commands)
- `ZOHO_INVENTORY_ORG_ID` - Default for `--org` flag (Inventory commands)
- `ZOHO_EXPENSE_ORG_ID` - Default for `--org` flag (Expense commands)
- `ZOHO_MAIL_ORG_ID` - Default for `--org` flag (Mail org-level APIs)
- `ZOHO_MAIL_ACCOUNT_ID` - Default for `--account` flag (Mail account-level APIs)
- `ZOHO_INVOICE_ORG_ID` - Default for `--org` flag (Invoice commands)
- `ZOHO_BOOKS_ORG_ID` - Default for `--org` flag (Books commands)

Flag passed on CLI always overrides the env var. If neither is set, commands fail with a clear error.

## Conventions

- No comments in code unless asked
- JSON output to stdout by default, errors to stderr
- Exit codes: 0=success, 1=general error, 2=auth error, 3=not found, 4=validation error
- Typed envelope structs for API responses, raw map[string]any for record data
- Pass through raw Zoho API responses (thin wrapper, no data transformation)
- --help-all shows jq-friendly output schemas per command

## Key Zoho API Quirks (absorb internally)

- CRM v8 requires `fields` param on list/related/notes/attachments endpoints
- CRM v8 search-global uses `searchword` param (not `word`)
- CRM v8 tags add/remove use JSON body (not query params)
- CRM pagination uses page_token for >2000 records
- Projects pagination: has_next_page can be string "true" or bool true
- WorkDrive uses JSON:API content-type (application/vnd.api+json)
- WorkDrive copy has reversed semantics: POST to destination with source in body
- WorkDrive file status codes: 1=active, 51=trash, 61=delete
- Writer R3002: empty documents cannot be exported
- Download endpoint: use workdrive.zoho.com/api/v1/download/{id} (not download.zoho.com)
- Zoho rate-limits: 10 access token refreshes per refresh_token per 10 minutes
- Go net/http needs explicit Accept: */* header (WorkDrive returns 415 without it)
