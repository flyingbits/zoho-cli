package books

import (
	"github.com/urfave/cli/v3"
)

func salesReceiptsCmd() *cli.Command {
	return &cli.Command{
		Name:  "salesreceipts",
		Usage: "Sales receipt operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List sales receipts", Flags: listFlags(), Action: runList("salesreceipts")},
			{Name: "get", Usage: "Get a sales receipt", ArgsUsage: "<salesreceipt-id>", Action: runGet("salesreceipts")},
			{Name: "create", Usage: "Create a sales receipt", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("salesreceipts")},
			{Name: "update", Usage: "Update a sales receipt", ArgsUsage: "<salesreceipt-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("salesreceipts")},
			{Name: "delete", Usage: "Delete a sales receipt", ArgsUsage: "<salesreceipt-id>", Action: runDelete("salesreceipts")},
			{Name: "email", Usage: "Email sales receipt", ArgsUsage: "<salesreceipt-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("salesreceipts", "email")},
		},
	}
}
