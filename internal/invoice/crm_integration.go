package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func crmIntegrationCmd() *cli.Command {
	return &cli.Command{
		Name:  "crm-integration",
		Usage: "Zoho CRM integration",
		Commands: []*cli.Command{
			{Name: "import-customer-by-account", Usage: "Import customer using CRM account ID", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}, &cli.StringFlag{Name: "crm-account-id", Required: true, Usage: "CRM account ID"}}, Action: crmImportCustomerByAccount},
			{Name: "import-customer-by-contact", Usage: "Import customer using CRM contact ID", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}, &cli.StringFlag{Name: "crm-contact-id", Required: true, Usage: "CRM contact ID"}}, Action: crmImportCustomerByContact},
			{Name: "import-item", Usage: "Import item using CRM product ID", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}, &cli.StringFlag{Name: "crm-product-id", Required: true, Usage: "CRM product ID"}}, Action: crmImportItem},
		},
	}
}

func crmImportCustomerByAccount(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	var body any
	json.Unmarshal([]byte(cmd.String("json")), &body)
	params := mergeParams(orgID, map[string]string{"crm_account_id": cmd.String("crm-account-id")})
	raw, err := c.Request("POST", c.InvoiceBase+"/contacts/crmimport", &zohttp.RequestOpts{JSON: body, Params: params})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func crmImportCustomerByContact(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	var body any
	json.Unmarshal([]byte(cmd.String("json")), &body)
	params := mergeParams(orgID, map[string]string{"crm_contact_id": cmd.String("crm-contact-id")})
	raw, err := c.Request("POST", c.InvoiceBase+"/contacts/crmimport", &zohttp.RequestOpts{JSON: body, Params: params})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func crmImportItem(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	var body any
	json.Unmarshal([]byte(cmd.String("json")), &body)
	params := mergeParams(orgID, map[string]string{"crm_product_id": cmd.String("crm-product-id")})
	raw, err := c.Request("POST", c.InvoiceBase+"/items/crmimport", &zohttp.RequestOpts{JSON: body, Params: params})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
