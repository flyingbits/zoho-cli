package inventory

import (
	"context"
	"encoding/json"
	"os"

	"github.com/omin8tor/zoho-cli/internal"
	"github.com/omin8tor/zoho-cli/internal/auth"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
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
	for c := cmd; c != nil; c = c.Root() {
		if c.Name == "inventory" {
			if org := c.String("org"); org != "" {
				return org, nil
			}
			break
		}
	}
	if org := os.Getenv("ZOHO_INVENTORY_ORG_ID"); org != "" {
		return org, nil
	}
	return "", internal.NewValidationError("--org flag or ZOHO_INVENTORY_ORG_ID env var required")
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
	return c.Request(method, c.InventoryBase+path, opts)
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "inventory",
		Usage: "Zoho Inventory API v1 operations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "org", Usage: "Organization ID (or set ZOHO_INVENTORY_ORG_ID)"},
		},
		Commands: []*cli.Command{
			organizationsCmd(),
			contactsCmd(),
			contactPersonsCmd(),
			itemGroupsCmd(),
			itemsCmd(),
			compositeItemsCmd(),
			itemAdjustmentsCmd(),
			transferOrdersCmd(),
			salesOrdersCmd(),
			packagesCmd(),
			shipmentOrdersCmd(),
			invoicesCmd(),
			retainerInvoicesCmd(),
			customerPaymentsCmd(),
			salesReturnsCmd(),
			creditNotesCmd(),
			purchaseOrdersCmd(),
			purchaseReceivesCmd(),
			billsCmd(),
			vendorCreditsCmd(),
			locationsCmd(),
			priceListsCmd(),
			usersCmd(),
			taxesCmd(),
			currencyCmd(),
			reportingTagsCmd(),
		},
	}
}

func crudSubcommands(path string) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "list",
			Usage: "List",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "page", Usage: "Page"},
				&cli.StringFlag{Name: "per-page", Usage: "Per page"},
				&cli.StringFlag{Name: "sort-column", Usage: "Sort column"},
				&cli.StringFlag{Name: "sort-order", Usage: "asc or desc"},
				&cli.StringFlag{Name: "filter-by", Usage: "Filter"},
			},
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				params := mergeParams(orgID, nil)
				if v := cmd.String("page"); v != "" {
					params["page"] = v
				}
				if v := cmd.String("per-page"); v != "" {
					params["per_page"] = v
				}
				if v := cmd.String("sort-column"); v != "" {
					params["sort_column"] = v
				}
				if v := cmd.String("sort-order"); v != "" {
					params["sort_order"] = v
				}
				if v := cmd.String("filter-by"); v != "" {
					params["filter_by"] = v
				}
				raw, err := req(c, orgID, "GET", "/"+path, &zohttp.RequestOpts{Params: params})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "get",
			Usage:     "Get by ID",
			ArgsUsage: "<id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "GET", "/"+path+"/"+cmd.Args().First(), nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:  "create",
			Usage: "Create",
			Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
			Action: func(_ context.Context, cmd *cli.Command) error {
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
				raw, err := req(c, orgID, "POST", "/"+path, &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "update",
			Usage:     "Update by ID",
			ArgsUsage: "<id>",
			Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
			Action: func(_ context.Context, cmd *cli.Command) error {
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
				raw, err := req(c, orgID, "PUT", "/"+path+"/"+cmd.Args().First(), &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "delete",
			Usage:     "Delete by ID",
			ArgsUsage: "<id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "DELETE", "/"+path+"/"+cmd.Args().First(), nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
}

