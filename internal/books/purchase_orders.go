package books

import (
	"github.com/urfave/cli/v3"
)

func purchaseOrdersCmd() *cli.Command {
	return &cli.Command{
		Name:  "purchaseorders",
		Usage: "Purchase order operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List purchase orders", Flags: listFlags(), Action: runList("purchaseorders")},
			{Name: "get", Usage: "Get a purchase order", ArgsUsage: "<purchaseorder-id>", Action: runGet("purchaseorders")},
			{Name: "create", Usage: "Create a purchase order", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("purchaseorders")},
			{Name: "update", Usage: "Update a purchase order", ArgsUsage: "<purchaseorder-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("purchaseorders")},
			{Name: "delete", Usage: "Delete a purchase order", ArgsUsage: "<purchaseorder-id>", Action: runDelete("purchaseorders")},
			{Name: "mark-open", Usage: "Mark as open", ArgsUsage: "<purchaseorder-id>", Action: runPost("purchaseorders", "status/open")},
			{Name: "mark-billed", Usage: "Mark as billed", ArgsUsage: "<purchaseorder-id>", Action: runPost("purchaseorders", "status/billed")},
			{Name: "cancel", Usage: "Cancel purchase order", ArgsUsage: "<purchaseorder-id>", Action: runPost("purchaseorders", "status/cancelled")},
			{Name: "submit", Usage: "Submit for approval", ArgsUsage: "<purchaseorder-id>", Action: runPost("purchaseorders", "submit")},
			{Name: "approve", Usage: "Approve purchase order", ArgsUsage: "<purchaseorder-id>", Action: runPost("purchaseorders", "approve")},
			{Name: "reject", Usage: "Reject purchase order", ArgsUsage: "<purchaseorder-id>", Action: runPost("purchaseorders", "reject")},
			{Name: "email", Usage: "Email purchase order", ArgsUsage: "<purchaseorder-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("purchaseorders", "email")},
			{Name: "list-comments", Usage: "List comments", ArgsUsage: "<purchaseorder-id>", Action: runGetSub("purchaseorders", "comments")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<purchaseorder-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("purchaseorders", "comments")},
		},
	}
}
