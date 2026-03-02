package mail

import (
	"os"

	"github.com/omin8tor/zoho-cli/internal"
	"github.com/omin8tor/zoho-cli/internal/auth"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

const orgRequiredMsg = "--org or ZOHO_MAIL_ORG_ID required"
const accountRequiredMsg = "--account or ZOHO_MAIL_ACCOUNT_ID required"

func getClient() (*zohttp.Client, error) {
	config, err := auth.ResolveAuth()
	if err != nil {
		return nil, err
	}
	return zohttp.NewClient(config)
}

func resolveOrg(cmd *cli.Command) (string, error) {
	org := cmd.String("org")
	if org == "" {
		org = os.Getenv("ZOHO_MAIL_ORG_ID")
	}
	if org == "" {
		return "", nil
	}
	return org, nil
}

func resolveAccount(cmd *cli.Command) (string, error) {
	acc := cmd.String("account")
	if acc == "" {
		acc = os.Getenv("ZOHO_MAIL_ACCOUNT_ID")
	}
	return acc, nil
}

func mustOrg(cmd *cli.Command) (string, error) {
	org, _ := resolveOrg(cmd)
	if org == "" {
		return "", internal.NewValidationError(orgRequiredMsg)
	}
	return org, nil
}

func mustAccount(cmd *cli.Command) (string, error) {
	acc, _ := resolveAccount(cmd)
	if acc == "" {
		return "", internal.NewValidationError(accountRequiredMsg)
	}
	return acc, nil
}

func request(c *zohttp.Client, method, path string, opts *zohttp.RequestOpts) error {
	raw, err := c.Request(method, c.MailBase+path, opts)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "mail",
		Usage: "Zoho Mail operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Usage: "Organization ID (zoid); or set ZOHO_MAIL_ORG_ID"},
			&cli.StringFlag{Name: "account", Usage: "Account ID for user-level APIs; or set ZOHO_MAIL_ACCOUNT_ID"},
		},
		Commands: []*cli.Command{
			organizationCmd(),
			domainsCmd(),
			groupsCmd(),
			usersCmd(),
			policyCmd(),
			accountsCmd(),
			foldersCmd(),
			labelsCmd(),
			messagesCmd(),
			signaturesCmd(),
			threadsCmd(),
			tasksCmd(),
			bookmarksCmd(),
			notesCmd(),
			logsCmd(),
		},
	}
}
