package books

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func crmIntegrationCmd() *cli.Command {
	return &cli.Command{
		Name:  "crmintegration",
		Usage: "Zoho CRM integration (import customer/vendor/item)",
		Commands: []*cli.Command{
			{
				Name:      "import-customer-by-account",
				Usage:     "Import a customer using the CRM account ID",
				ArgsUsage: "<crm-account-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Params: mergeParams(orgID, map[string]string{"crm_account_id": cmd.Args().First()})}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
						opts.JSON = body
					}
					raw, err := req(c, orgID, "POST", "/contacts/crmimport", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "import-customer-by-contact",
				Usage:     "Import a customer using CRM contact ID",
				ArgsUsage: "<crm-contact-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Params: mergeParams(orgID, map[string]string{"crm_contact_id": cmd.Args().First()})}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
						opts.JSON = body
					}
					raw, err := req(c, orgID, "POST", "/contacts/crmimport", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "import-vendor",
				Usage:     "Import a vendor using the CRM vendor ID",
				ArgsUsage: "<crm-vendor-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Params: mergeParams(orgID, map[string]string{"crm_vendor_id": cmd.Args().First()})}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
						opts.JSON = body
					}
					raw, err := req(c, orgID, "POST", "/vendors/crmimport", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "import-item",
				Usage:     "Import an item using the CRM product ID",
				ArgsUsage: "<crm-product-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{Params: mergeParams(orgID, map[string]string{"crm_product_id": cmd.Args().First()})}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
						opts.JSON = body
					}
					raw, err := req(c, orgID, "POST", "/items/crmimport", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
