package books

import (
	"github.com/urfave/cli/v3"
)

func salesOrdersCmd() *cli.Command {
	return &cli.Command{
		Name:  "salesorders",
		Usage: "Sales order operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List sales orders", Flags: listFlags(), Action: runList("salesorders")},
			{Name: "get", Usage: "Get a sales order", ArgsUsage: "<salesorder-id>", Action: runGet("salesorders")},
			{Name: "create", Usage: "Create a sales order", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("salesorders")},
			{Name: "update", Usage: "Update a sales order", ArgsUsage: "<salesorder-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("salesorders")},
			{Name: "delete", Usage: "Delete a sales order", ArgsUsage: "<salesorder-id>", Action: runDelete("salesorders")},
			{Name: "mark-open", Usage: "Mark as open", ArgsUsage: "<salesorder-id>", Action: runPost("salesorders", "status/open")},
			{Name: "mark-void", Usage: "Mark as void", ArgsUsage: "<salesorder-id>", Action: runPost("salesorders", "status/void")},
			{Name: "email", Usage: "Email sales order", ArgsUsage: "<salesorder-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("salesorders", "email")},
			{Name: "submit", Usage: "Submit for approval", ArgsUsage: "<salesorder-id>", Action: runPost("salesorders", "submit")},
			{Name: "approve", Usage: "Approve sales order", ArgsUsage: "<salesorder-id>", Action: runPost("salesorders", "approve")},
			{Name: "list-comments", Usage: "List comments", ArgsUsage: "<salesorder-id>", Action: runGetSub("salesorders", "comments")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<salesorder-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("salesorders", "comments")},
		},
	}
}
