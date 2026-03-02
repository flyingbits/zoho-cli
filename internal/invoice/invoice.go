package invoice

import (
	"encoding/json"
	"os"

	"github.com/omin8tor/zoho-cli/internal"
	"github.com/omin8tor/zoho-cli/internal/auth"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func getClient() (*zohttp.Client, error) {
	config, err := auth.ResolveAuth()
	if err != nil {
		return nil, err
	}
	return zohttp.NewClient(config)
}

func resolveOrgID(cmd *cli.Command) (string, error) {
	for c := cmd; c != nil; c = c.Parent() {
		if c.Name == "invoice" {
			if org := c.String("org"); org != "" {
				return org, nil
			}
			break
		}
	}
	if org := os.Getenv("ZOHO_INVOICE_ORG_ID"); org != "" {
		return org, nil
	}
	return "", internal.NewValidationError("--org flag or ZOHO_INVOICE_ORG_ID env var required")
}

func orgParam(orgID string) map[string]string {
	return map[string]string{"organization_id": orgID}
}

func mergeParams(orgID string, extra map[string]string) map[string]string {
	m := map[string]string{"organization_id": orgID}
	for k, v := range extra {
		m[k] = v
	}
	return m
}

func req(c *zohttp.Client, orgID, method, path string, opts *zohttp.RequestOpts) (json.RawMessage, error) {
	if opts == nil {
		opts = &zohttp.RequestOpts{}
	}
	if opts.Params == nil {
		opts.Params = orgParam(orgID)
	} else {
		opts.Params["organization_id"] = orgID
	}
	return c.Request(method, c.InvoiceBase+path, opts)
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "invoice",
		Usage: "Zoho Invoice API v3 operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Usage: "Organization ID (or set ZOHO_INVOICE_ORG_ID)"},
		},
		Commands: []*cli.Command{
			organizationsCmd(),
			itemsCmd(),
			priceListsCmd(),
			contactsCmd(),
			contactPersonsCmd(),
			estimatesCmd(),
			invoicesCmd(),
			recurringInvoicesCmd(),
			customerPaymentsCmd(),
			retainerInvoicesCmd(),
			creditNotesCmd(),
			expensesCmd(),
			recurringExpensesCmd(),
			projectsCmd(),
			tasksCmd(),
			timeEntriesCmd(),
			usersCmd(),
			taxesCmd(),
			expenseCategoryCmd(),
			currencyCmd(),
			crmIntegrationCmd(),
		},
	}
}
