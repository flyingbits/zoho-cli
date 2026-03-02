package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func itemsCmd() *cli.Command {
	return &cli.Command{
		Name:  "items",
		Usage: "Item operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List items",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Page"},
					&cli.StringFlag{Name: "per-page", Usage: "Per page"},
					&cli.StringFlag{Name: "sort-column", Usage: "Sort column"},
					&cli.StringFlag{Name: "sort-order", Usage: "asc or desc"},
					&cli.StringFlag{Name: "name", Usage: "Filter by name"},
					&cli.StringFlag{Name: "description", Usage: "Filter by description"},
					&cli.StringFlag{Name: "item-id", Usage: "Filter by item ID"},
					&cli.StringFlag{Name: "sku", Usage: "Filter by SKU"},
					&cli.StringFlag{Name: "type", Usage: "Filter by type"},
					&cli.StringFlag{Name: "status", Usage: "active or inactive"},
				},
				Action: runList("items"),
			},
			{
				Name:      "get",
				Usage:     "Retrieve an item",
				ArgsUsage: "<item-id>",
				Action:    runGet("items"),
			},
			{
				Name:  "create",
				Usage: "Create an item",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: runCreate("items"),
			},
			{
				Name:      "update",
				Usage:     "Update an item",
				ArgsUsage: "<item-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action:    runUpdate("items"),
			},
			{
				Name:      "delete",
				Usage:     "Delete an item",
				ArgsUsage: "<item-id>",
				Action:    runDelete("items"),
			},
			{
				Name:  "bulk-get",
				Usage: "Bulk fetch item details",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "item_ids", Required: true, Usage: "Comma-separated item IDs"},
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
					params := mergeParams(orgID, map[string]string{"item_ids": cmd.String("item_ids")})
					raw, err := c.Request("GET", c.InvoiceBase+"/items", &zohttp.RequestOpts{Params: params})
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
					raw, err := req(c, orgID, "PUT", "/items/"+cmd.Args().First(), &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "mark-active",
				Usage:     "Mark item as active",
				ArgsUsage: "<item-id>",
				Action:    runPost("items", "active"),
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark item as inactive",
				ArgsUsage: "<item-id>",
				Action:    runPost("items", "inactive"),
			},
		},
	}
}
