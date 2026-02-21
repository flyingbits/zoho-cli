# zoho-cli

CLI for Zoho's REST APIs. Covers CRM, Projects, WorkDrive, Writer, and Cliq — 118 commands, single binary, JSON to stdout.

No other tool does this. Zoho's own CLIs (ZET, Catalyst) are for building extensions, not talking to the API. Their Python SDK covers CRM only and requires MySQL for token storage.

## Getting started

You'll need two things: the binary and a Zoho API app. Takes about five minutes.

### Install

From source (Go 1.22+):

```bash
go install github.com/omin8tor/zoho-cli/cmd/zoho@latest
```

Or grab a prebuilt binary from [releases](https://github.com/omin8tor/zoho-cli/releases) — unpack it, put it on your PATH.

### Set up a Zoho API app

Go to [api-console.zoho.com](https://api-console.zoho.com/) and create a "Self Client" app. You'll get a client ID and secret.

### Authenticate

**If you're a human at a terminal**, use device flow:

```bash
zoho auth login --client-id YOUR_ID --client-secret YOUR_SECRET
```

This opens your browser, you approve the scopes, and you're done. Tokens are saved to `~/.config/zoho-cli/config.json` and auto-refresh.

**If you're a script or agent**, use env vars:

```bash
export ZOHO_CLIENT_ID=1000.ABC123
export ZOHO_CLIENT_SECRET=xyz789
export ZOHO_REFRESH_TOKEN=1000.refresh_token_here
export ZOHO_DC=com
```

(You can also use `zoho auth self-client` to exchange a code from the API Console.)

Check that it works:

```bash
zoho auth status
```

## Tutorial: your first queries

Everything below assumes you've authenticated. Output is always JSON to stdout, errors to stderr. Pipe into `jq` for filtering.

### CRM: look up contacts

List contacts, pulling specific fields:

```bash
zoho crm records list Contacts --fields "Full_Name,Email,Phone"
```

Get one contact by ID:

```bash
zoho crm records get Contacts 5551234000000012345 --fields "Full_Name,Email"
```

Search by criteria — find closed-won deals:

```bash
zoho crm records search Deals --criteria "(Stage:equals:Closed Won)" --fields "Deal_Name,Amount"
```

Create a lead:

```bash
zoho crm records create Leads --json '{"Last_Name":"Ochoa","Company":"Acme"}'
```

CRM v8 requires the `--fields` param on most read endpoints. If you forget it, the API returns empty records. That's Zoho, not us.

### CRM: search across everything

```bash
zoho crm search-global --searchword "Ochoa" --fields "Full_Name,Email"
```

### CRM: COQL queries

If criteria-based search isn't flexible enough, use COQL (Zoho's SQL-like query language):

```bash
zoho crm coql --query "SELECT Full_Name, Email FROM Contacts WHERE Email LIKE '%@acme.com' LIMIT 10"
```

This needs the `ZohoCRM.coql.READ` scope — it's separate from the general CRM scopes.

### Projects: find your tasks

You need a portal ID. List your portals first:

```bash
zoho projects core list
```

Grab the portal ID from the output, then:

```bash
zoho projects tasks my --portal 12345
```

List all tasks in a project:

```bash
zoho projects tasks list --portal 12345 --project 67890
```

Filter open tasks with jq:

```bash
zoho projects tasks my --portal 12345 | jq '[.[] | select(.status.name == "Open")]'
```

### WorkDrive: navigate and download

Find your team:

```bash
zoho drive teams me
```

List top-level folders:

```bash
zoho drive folders list --team TEAM_ID
```

List files in a folder:

```bash
zoho drive files list --folder FOLDER_ID
```

Download a file:

```bash
zoho drive download FILE_ID --output ./report.pdf
```

Upload a file:

```bash
zoho drive upload ./quarterly.xlsx --folder FOLDER_ID
```

### Writer: work with documents

Get document details:

```bash
zoho writer details DOC_ID
```

Merge data into a template and export as PDF:

```bash
zoho writer merge DOC_ID --json '{"name":"Alice","date":"2025-01-15"}' --format pdf --output ./letter.pdf
```

Get doc IDs from WorkDrive — Writer has no "list documents" endpoint.

### Cliq: send messages

List channels:

```bash
zoho cliq channels list
```

Send a DM:

```bash
zoho cliq buddies message someone@company.com --text "quarterly report is ready"
```

## Piping and composition

The whole point is composability. Everything is JSON, so chain with `jq`, `xargs`, whatever:

```bash
# Get all emails from contacts
zoho crm records list Contacts --fields "Email" | jq -r '.[].Email'

# Download every file in a folder
zoho drive files list --folder FOLDER_ID | jq -r '.[].id' | xargs -I{} zoho drive download {} --output ./downloads/

# Find overdue tasks
zoho projects tasks my --portal 12345 | jq '[.[] | select(.end_date < "2025-01-01")]'
```

## Data centers

Zoho runs in 9 data centers. Set via `ZOHO_DC` env var or `--dc` flag on auth commands:

`com` (US, default) · `eu` · `in` · `com.au` · `jp` · `ca` · `sa` · `uk` · `com.cn`

## Supported products

| Product | Status | Commands |
|---------|--------|----------|
| **CRM** | Supported | 29 |
| **Projects** | Supported | 39 |
| **WorkDrive** | Supported | 26 |
| **Writer** | Supported | 7 |
| **Cliq** | Supported | 12 |
| Desk | Planned | — |
| Books | Planned | — |
| People | Planned | — |
| Recruit | Planned | — |
| Analytics | Planned | — |
| Sign | Planned | — |
| Campaigns | Planned | — |
| Mail | Planned | — |
| Calendar | Planned | — |
| Sheet | Planned | — |
| Show | Planned | — |
| Inventory | Planned | — |
| Invoice | Planned | — |
| Expense | Planned | — |
| Billing | Planned | — |
| Forms | Planned | — |
| SalesIQ | Planned | — |
| Bookings | Planned | — |
| Social | Planned | — |
| Survey | Planned | — |
| Meeting | Planned | — |
| Connect | Planned | — |
| Flow | Planned | — |
| Creator | Planned | — |
| Sprints | Planned | — |
| BugTracker | Planned | — |
| Bigin | Planned | — |
| Voice | Planned | — |
| Commerce | Planned | — |
| Backstage | Planned | — |
| Marketing Automation | Planned | — |
| FSM | Planned | — |
| Assist | Planned | — |
| Directory | Planned | — |
| Shifts | Planned | — |
| Contracts | Planned | — |
| Practice | Planned | — |
| Checkout | Planned | — |
| Lens | Planned | — |
| Learn | Planned | — |
| ZeptoMail | Planned | — |
| Notebook | Planned | — |
| TeamInbox | Planned | — |
| Office Integrator | Planned | — |
| ToDo | Planned | — |
| PDF Editor | Planned | — |
| IoT | Planned | — |
| DataPrep | Planned | — |
| Apptics | Planned | — |
| Vault | Planned | — |
| Catalyst | Planned | — |
| Webinar | Planned | — |
| PageSense | Planned | — |
| LandingPage | Planned | — |
| CommunitySpaces | Planned | — |
| Thrive | Planned | — |
| Sites | Planned | — |
| RouteIQ | Planned | — |
| Workerly | Planned | — |
| Solo | Planned | — |
| Procurement | Planned | — |

Want a product prioritized? [Open an issue](https://github.com/omin8tor/zoho-cli/issues).

Run `zoho --help-all` for the full command reference with every flag.

## Agent Skill

This repo ships as an [Agent Skill](https://agentskills.io/) so LLM agents can discover and use it. The skill definition is in [`SKILL.md`](./SKILL.md) with detailed references in [`references/`](./references/).

## Development

```bash
go build -o zoho ./cmd/zoho/
go test ./...
go vet ./...
```

Or with [mise](https://mise.jdx.dev/):

```bash
mise run build    # build binary
mise run test     # unit tests
mise run lint     # go vet
```

## License

GPL-3.0
