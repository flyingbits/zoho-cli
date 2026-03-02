package inventory

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func itemsCmd() *cli.Command {
	base := crudSubcommands("items")
	extra := []*cli.Command{
		{
			Name:  "bulk-get",
			Usage: "Bulk fetch item details",
			Flags: []cli.Flag{&cli.StringFlag{Name: "item-ids", Required: true, Usage: "Comma-separated item IDs"}},
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				params := mergeParams(orgID, map[string]string{"item_ids": cmd.String("item-ids")})
				raw, err := req(c, orgID, "GET", "/items", &zohttp.RequestOpts{Params: params})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "update-custom-field",
			Usage:     "Update custom field in existing items",
			ArgsUsage: "<item-id>",
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
				raw, err := req(c, orgID, "PUT", "/items/"+cmd.Args().First()+"/customfield", &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "delete-image",
			Usage:     "Delete item image",
			ArgsUsage: "<item-id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "DELETE", "/items/"+cmd.Args().First()+"/image", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-active",
			Usage:     "Mark item as active",
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
				raw, err := req(c, orgID, "POST", "/items/"+cmd.Args().First()+"/active", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-inactive",
			Usage:     "Mark item as inactive",
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
				raw, err := req(c, orgID, "POST", "/items/"+cmd.Args().First()+"/inactive", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "items",
		Usage:    "Item operations",
		Commands: append(base, extra...),
	}
}
